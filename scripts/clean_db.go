package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
	"video-platform/config"
	"video-platform/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	all      bool
	days     int
	dryRun   bool
	testOnly bool
)

func init() {
	flag.BoolVar(&all, "all", false, "清理所有数据")
	flag.IntVar(&days, "days", 30, "清理多少天前的数据")
	flag.BoolVar(&dryRun, "dry-run", false, "试运行模式，不实际删除数据")
	flag.BoolVar(&testOnly, "test-only", false, "只清理测试数据")
	flag.Parse()
}

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("配置初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库连接
	if err := database.InitMongoDB(); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer database.CloseMongoDB()

	ctx := context.Background()
	filter := bson.M{}

	// 构建查询条件
	if !all {
		if testOnly {
			// 只清理测试数据
			filter["user_id"] = bson.M{"$regex": "^test_"}
		} else {
			// 清理指定天数前的数据
			deadline := time.Now().AddDate(0, 0, -days)
			filter["created_at"] = bson.M{"$lt": deadline}
		}
	}

	// 获取要删除的数据
	collection := database.GetCollection("videos")
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Printf("查询数据失败: %v\n", err)
		os.Exit(1)
	}

	// 打印将要删除的数据信息
	fmt.Printf("将要删除 %d 条数据\n", count)
	if count == 0 {
		fmt.Println("没有需要删除的数据")
		return
	}

	// 如果是试运行模式，到这里就结束
	if dryRun {
		fmt.Println("试运行模式，不实际删除数据")
		return
	}

	// 确认是否继续
	if !confirm() {
		fmt.Println("操作已取消")
		return
	}

	// 执行删除
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Printf("删除数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功删除 %d 条数据\n", result.DeletedCount)
}

// confirm 确认是否继续
func confirm() bool {
	fmt.Print("确认要删除这些数据吗？(y/N): ")
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
