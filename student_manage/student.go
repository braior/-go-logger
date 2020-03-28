package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Student struct
type Student struct {
	ID     int
	Name   string
	Gender string
	Age    int8
	Score  int8
}

// StudentMgr struct
type StudentMgr struct {
	allStudent map[int]*Student
}

func newStudent(id int, name, gender string, age, score int8) *Student {
	return &Student{
		ID:     id,
		Name:   name,
		Age:    age,
		Gender: gender,
		Score:  score,
	}
}

func newStudentMgr() *StudentMgr {
	return &StudentMgr{
		allStudent: make(map[int]*Student, 50),
	}
}

func checkInput(sm *StudentMgr, id int, value interface{}) bool {
	// switch v := value.(type) {
	// case string:
	// 	fmt.Println(v)
	// 	return false
	// }

	switch id {
	case 1001:
		fmt.Println(value)
		return false
	default:
		return true
	}
}

func getInputAdd(sm *StudentMgr) (int, interface{}) {
	var (
		id     int
		name   string
		gender string
		age    int8
		score  int8
	)

	stdin := bufio.NewReader(os.Stdin)
	fmt.Print("请输入学员ID：")
	_, err := fmt.Scanf("%d\n", &id)
	_, isSave := sm.allStudent[id]
	if err != nil {
		stdin.ReadString('\n')
		return 1001, "输入的学员id错误！"
	} else if isSave {
		return 1001, "学生ID已存在！"
	}

	fmt.Print("请输入姓名：")
	_, err = fmt.Scanf("%s\n", &name)
	if err != nil {
		stdin.ReadString('\n')
		return 1001, "输入的学员姓名错误！"
	}

	fmt.Print("请输入学员性别[男|女]：")
	fmt.Scanf("%s\n", &gender)
	if !(strings.EqualFold(gender, "男") || strings.EqualFold(gender, "女")) {
		return 1001, "输入的学员性别错误！"
	}

	fmt.Print("请输入学生年龄：")
	_, err = fmt.Scanf("%d\n", &age)
	if err != nil {
		stdin.ReadString('\n')
		return 1001, "输入的学员年龄错误！"
	}

	fmt.Print("请输入学生分数：")
	_, err = fmt.Scanf("%d\n", &score)
	if err != nil {
		stdin.ReadString('\n')
		return 1001, "输入的学员分数错误！"
	}

	stu := newStudent(id, name, gender, age, score)
	return id, stu
}

func (sm *StudentMgr) modiftyStudent() {
	var id int
	stdin := bufio.NewReader(os.Stdin)
	fmt.Print("输入要修改学生的ID：")
	_, err := fmt.Scanf("%d\n", &id)
	_, isSave := sm.allStudent[id]
	if err != nil {
		stdin.ReadString('\n')
		fmt.Println("输入的学生id错误！")
		return
	} else if !isSave {
		fmt.Println("要修改的学生ID不存在！")
		return
	}

	fmt.Print("请输入要修改的属性[1.姓名 2.性别 3.年龄 4.分数]:")
	var choice int
	_, err = fmt.Scanf("%d\n", &choice)
	if err != nil {
		stdin.ReadString('\n')
		fmt.Println("输入的选项错误！")
		return
	}

	switch choice {
	case 1:
		var newName string
		fmt.Print("please input new name:")
		_, err := fmt.Scanf("%s\n", &newName)
		if err != nil {
			stdin.ReadString('\n')
			fmt.Println("输入的姓名错误！")
		} else {
			sm.allStudent[id].Name = newName
			fmt.Println("已成功修改学生姓名！")
		}

	case 2:
		var newGender string
		fmt.Print("please input new gender:")
		fmt.Scanf("%s\n", &newGender)
		if !(strings.EqualFold(newGender, "男") || strings.EqualFold(newGender, "女")) {
			fmt.Println("输入的学员性别错误！")
		} else {
			sm.allStudent[id].Gender = newGender
			fmt.Println("已成功修改学生性别！")
		}

	case 3:
		var newAge int8
		fmt.Print("please input new age:")
		_, err := fmt.Scanf("%d\n", &newAge)
		if err != nil {
			fmt.Println("输入的学员年龄错误！")
			stdin.ReadString('\n')
		} else {
			sm.allStudent[id].Age = newAge
			fmt.Println("已成功修改学生年龄！")
		}

	case 4:
		var newScore int8
		fmt.Print("please input new score:")
		_, err := fmt.Scanf("%d\n", &newScore)
		if err != nil {
			fmt.Println("输入的分数错误！")
			stdin.ReadString('\n')
		} else {
			sm.allStudent[id].Score = newScore
			fmt.Println("已成功修改学生分数！")
		}
	default:
		fmt.Println("输入的选项错误！")
	}
}

func (sm *StudentMgr) addStudent(id int, newStu *Student) {
	sm.allStudent[id] = newStu
}

func (sm *StudentMgr) showStudent() {
	fmt.Printf("\t%s\t%s\t%s\t%s\t%s\n", "ID", "姓名", "性别", "年龄", "分数")
	sortID := make([]int, 0)
	for k := range sm.allStudent {
		sortID = append(sortID, k)
	}
	sort.Ints(sortID)
	for _, k := range sortID {
		fmt.Printf("\t%d\t%s\t%s\t%d\t%d\n", sm.allStudent[k].ID, sm.allStudent[k].Name,
			sm.allStudent[k].Gender, sm.allStudent[k].Age, sm.allStudent[k].Score)
	}
}

func (sm *StudentMgr) deleteStudent() {
	stdin := bufio.NewReader(os.Stdin)

	var id int
	fmt.Print("请输入要删除学生的ID:")
	_, err := fmt.Scanf("%d\n", &id)
	if err != nil {
		fmt.Println("输入错误！")
		stdin.ReadString('\n')
		return
	}

	_, isSave := sm.allStudent[id]
	if !isSave {
		fmt.Println("不存在的ID!")
		return
	}
	delete(sm.allStudent, id)
	fmt.Println("删除成功！")
}
