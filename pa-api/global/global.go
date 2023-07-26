package global

var instanceId int = 0

func IncrementInstanceId() {
	instanceId = instanceId + 1
}

func GetInstanceId() int {
	return instanceId
}
