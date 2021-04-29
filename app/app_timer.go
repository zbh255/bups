package app

/*
	根据应用程序的定时器状态执行的任务
*/

func TimerTask() {
	_ = CreateSqlFileZip()
	BackUpForFile()
	ReadZipFile()
	BackUpForDb()
}
