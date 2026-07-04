package main

import (
	"awesomeProject/ent"
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

func Example() {
	// 使用内存中SQLite数据库来创建 ent.Client
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// 运行自动迁移工具来创建所有Schema资源
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// 输出:

	// ...
	task1, err := client.Todo.Create().Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	fmt.Println(task1)
	// Output:
	// Todo(id=1)
}
