package script

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
	dataType string
	help     bool
)

func init() {
	flag.BoolVar(&all, "all", false, "清理所有数据")
	flag.IntVar(&days, "days", 30, "清理多少天前的数据")
	flag.BoolVar(&dryRun, "dry-run", false, "试运行模式，不实际删除数据")
	flag.BoolVar(&testOnly, "test-only", false, "只清理测试数据")
	flag.StringVar(&dataType, "type", "all", "要清理的数据类型(all|videos|marks|annotations|notes)")
	flag.BoolVar(&help, "h", false, "显示帮助信息")
	flag.Usage = usage
}

func usage() {
	fmt.Println("数据库清理工具")
	fmt.Println("\n用法:")
	fmt.Println("  clean_db [选项]")
	fmt.Println("\n选项:")
	fmt.Println("  -all         清理所有数据")
	fmt.Println("  -days n      清理n天前的数据（默认30天）")
	fmt.Println("  -dry-run     试运行模式，不实际删除数据")
	fmt.Println("  -test-only   只清理测试数据")
	fmt.Println("  -type string 要清理的数据类型：")
	fmt.Println("               all: 所有数据（默认）")
	fmt.Println("               videos: 视频数据")
	fmt.Println("               marks: 标记数据")
	fmt.Println("               annotations: 注释数据")
	fmt.Println("               notes: 笔记数据")
	fmt.Println("  -h          显示帮助信息")
	fmt.Println("\n示例:")
	fmt.Println("  清理30天前的所有数据:")
	fmt.Println("    clean_db -days 30")
	fmt.Println("  清理所有测试数据:")
	fmt.Println("    clean_db -test-only")
	fmt.Println("  试运行模式查看将要删除的数据:")
	fmt.Println("    clean_db -dry-run")
	fmt.Println("  只清理视频数据:")
	fmt.Println("    clean_db -type videos")
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	// 验证数据类型参数
	validTypes := map[string]bool{
		"all":         true,
		"videos":      true,
		"marks":       true,
		"annotations": true,
		"notes":       true,
	}
	if !validTypes[dataType] {
		fmt.Printf("无效的数据类型: %s\n", dataType)
		flag.Usage()
		os.Exit(1)
	}

	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("配置初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库连接
	if err := database.InitMongoDB(context.Background(), config.GlobalConfig.MongoDB, false); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer database.CloseMongoDB()

	ctx := context.Background()
	collections := getCollections()

	for _, col := range collections {
		if err := cleanCollection(ctx, col); err != nil {
			fmt.Printf("清理集合 %s 失败: %v\n", col, err)
		}
	}
}

func getCollections() []string {
	switch dataType {
	case "all":
		return []string{"videos", "marks", "annotations", "notes"}
	case "videos":
		return []string{"videos"}
	case "marks":
		return []string{"marks"}
	case "annotations":
		return []string{"annotations"}
	case "notes":
		return []string{"notes"}
	default:
		return []string{}
	}
}

func cleanCollection(ctx context.Context, colName string) error {
	collection := database.GetCollection(colName)
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

	// 获取要删除的数据数量
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("查询数据失败: %v", err)
	}

	// 打印将要删除的数据信息
	fmt.Printf("\n集合 %s 将删除 %d 条数据\n", colName, count)
	if count == 0 {
		return nil
	}

	// 如果是试运行模式，到这里就结束
	if dryRun {
		return nil
	}

	// 确认是否继续
	if !confirm(colName) {
		fmt.Printf("跳过清理集合 %s\n", colName)
		return nil
	}

	// 执行删除
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除数据失败: %v", err)
	}

	fmt.Printf("成功从集合 %s 删除 %d 条数据\n", colName, result.DeletedCount)
	return nil
}

// confirm 确认是否继续
func confirm(colName string) bool {
	fmt.Printf("确认要删除集合 %s 中的这些数据吗？(y/N): ", colName)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
