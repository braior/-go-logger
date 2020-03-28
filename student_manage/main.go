package main

import (
	"fmt"
	"os"
	"time"
)

func menu() {
	fmt.Println("欢迎来到学院管理系统")
	fmt.Println("1.添加学员信息")
	fmt.Println("2.查看学员信息")
	fmt.Println("3.修改学员信息")
	fmt.Println("4.删除学员信息")
	fmt.Println("5.退出程序")
}

func main() {
	sm := newStudentMgr()
	for {
		menu()
		fmt.Print("请选择要执行的操作:")
		var do int8
		fmt.Scanf("%d\n", &do)

		switch do {
		case 1:
			id, stu := getInputAdd(sm)
			//stuTmp := check(id, stu)
			//sm.addStudent(id, stuTmp)
			if checkInput(sm, id, stu) {
				val, _ := stu.(*Student)
				sm.addStudent(id, val)
				fmt.Println("添加学员成功！")
			}
			time.Sleep(time.Duration(2) * time.Second)
		case 2:
			sm.showStudent()
			time.Sleep(time.Duration(2) * time.Second)
		case 3:
			sm.modiftyStudent()
			time.Sleep(time.Duration(2) * time.Second)
		case 4:
			sm.deleteStudent()
			time.Sleep(time.Duration(2) * time.Second)
		case 5:
			os.Exit(0)
		default:
			fmt.Println("输入错误！无效的选择")
		}
	}
}
