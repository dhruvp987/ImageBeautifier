/*
 * Authors: Dhruv Patel and Ayush Sharma
 * File: workerpool.go
 * Description:
 *   A worker pool to manage concurrent tasks throughout
 *   the program.
 */

package main

import (
    "sync"
)

/*
 * The procedure that each worker will be running. Each worker
 * will wait for a task, run it, and repeat until no more tasks
 * are incoming.
 */
func work(tasksCh <-chan func(), wg *sync.WaitGroup) {
    defer wg.Done()
    for task := range tasksCh {
        task()
    }
}

/*
 * A Worker Pool can accept incoming tasks and run them concurrently.
 */
type WorkerPool struct {
    NumWorkers int
    tasksCh chan func();
    wg *sync.WaitGroup;
}

/*
 * Construct a new Worker Pool. It will not be started yet.
 */
func WpNew(numWorkers int, chBufSize int) WorkerPool {
    var pool WorkerPool
    var wg sync.WaitGroup

    pool.NumWorkers = numWorkers
    pool.tasksCh = make(chan func(), chBufSize)
    pool.wg = &wg

    return pool
}

/*
 * Start a pool of numWorkers workers.
 */
func (p *WorkerPool) Start() {
    p.wg.Add(p.NumWorkers)
    for i := 0; i < p.NumWorkers; i++ {
        go work(p.tasksCh, p.wg)
    }
}

/*
 * Submit a task to the Worker Pool. This will wait until the task
 * has been submitted.
 */
func (p *WorkerPool) Submit(f func()) {
    p.tasksCh <- f
}

/*
 * Stop all workers in the pool. The function will wait for all tasks
 * to finish before returning.
 */
func (p *WorkerPool) WaitAndStop() {
    // Alert workers that there will be no more tasks sent.
    close(p.tasksCh)
    p.wg.Wait()
}
