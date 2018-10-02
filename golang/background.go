package golang



import (
	"fmt"
	"sync"
	"time"
	"log"
)

var (
	BackgroundTaskManager *BackgroundManager
)

// creates & starts workers
type BackgroundManager struct {
	numWorkers   int
	maxRetries   int
	stopInterval int64
	stopChannel  chan bool
	taskChannel  chan Task
	waitAllDone  *sync.WaitGroup
}

func NewBackgroundManager(numWorkers, maxRetries, maxBuffer int,
	stopInterval int64) *BackgroundManager {
	bgmanager := &BackgroundManager{
		numWorkers:   numWorkers,
		maxRetries:   maxRetries,
		stopInterval: stopInterval,
		stopChannel:  make(chan bool),
		taskChannel:  make(chan Task, maxBuffer),
		waitAllDone:  new(sync.WaitGroup),
	}
	bgmanager.Start()
	return bgmanager
}

// Put puts task for proccessing in background
func (mng *BackgroundManager) Put(execFunc func() error, description string) {
	select {
	case mng.taskChannel <- Task{
		execute:     execFunc,
		sleeper:     NewExponentialSleeper(time.Second, 10*time.Second),
		putTime:     time.Now(),
		description: description,
	}:
	default:
		log.Errorf("Task queue buffer is full while trying to append task %s, skipping", description)
	}
}

// Start all workers for proccesing tasks
func (mng *BackgroundManager) Start() {
	for i := 0; i < mng.numWorkers; i++ {
		mng.waitAllDone.Add(1)
		go mng.ProcessTasks()
	}
}

// Stop all workers of pool waiting until all finished
func (mng *BackgroundManager) Stop() {
	close(mng.stopChannel)
	mng.waitAllDone.Wait()
	if len(mng.taskChannel) != 0 {
		mng.flush()
	}
}

func (mng *BackgroundManager) ProcessTasks() {
	var stopped bool
	for !stopped {
		select {
		case task := <-mng.taskChannel:
			mng.processTask(task, true)
		case <-mng.stopChannel:
			stopped = true
			break
		}
	}
	mng.waitAllDone.Done()
}

func (mng *BackgroundManager) processTask(task Task, sleeping bool) {
	if task.failed && sleeping {
		task.sleeper.Decrease(time.Now().Sub(task.putTime))
		task.sleeper.Sleep()
	}
	err := task.execute()
	task.attempts += 1
	if err != nil {
		if task.attempts <= mng.maxRetries {
			task.failed = true
			task.putTime = time.Now()
			mng.taskChannel <- task // retry & append to end!
			log.Warning("BG WORKER. task %s failed: %s",
				task.String(), err.Error())
		} else {
			log.Error("BG WORKER. task %s failed execution: %s",
				task.String(), err.Error())
		}
	}
}

func (mng *BackgroundManager) flush() {
	// 25 % of stop interval on logging & 75 % on processing tasks
	interval := mng.stopInterval - mng.stopInterval/4

	timeout := time.NewTimer(time.Duration(interval) * time.Second)

	finished := false
	for !finished {
		select {
		case task := <-mng.taskChannel:
			mng.processTask(task, false)
			if len(mng.taskChannel) == 0 {
				finished = true
			}
		case <-timeout.C:
			finished = true
		}
	}

	if len(mng.taskChannel) == 0 {
		return
	}

	finished = false
	for !finished {
		select {
		case task := <-mng.taskChannel:
			log.Error("BG WORKER. No time for task: %s", task.String())
			if len(mng.taskChannel) == 0 {
				finished = true
			}
		}
	}
}

type Task struct {
	execute     func() error
	sleeper     *ExponentialSleeper
	attempts    int
	putTime     time.Time
	failed      bool
	description string
}

func (task *Task) String() string {
	return fmt.Sprintf("%s: retries - %d", task.description, task.attempts)
}
