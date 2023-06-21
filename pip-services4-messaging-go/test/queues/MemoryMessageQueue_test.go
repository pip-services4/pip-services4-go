package test_queues

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

func TestMemoryMessageQueue(t *testing.T) {
	queue := queues.NewMemoryMessageQueue("TestQueue")
	fixture := NewMessageQueueFixture(queue)

	queue.Open(context.TODO())
	defer queue.Close(context.TODO())
	queue.Clear(context.TODO())

	t.Run("MemoryMessageQueue:Send Receive Message", fixture.TestSendReceiveMessage)
	t.Run("MemoryMessageQueue:Receive Send Message", fixture.TestReceiveSendMessage)
	t.Run("MemoryMessageQueue:Receive And Complete Message", fixture.TestReceiveCompleteMessage)
	t.Run("MemoryMessageQueue:Receive And Abandon Message", fixture.TestReceiveAbandonMessage)
	t.Run("MemoryMessageQueue:Send Peek Message", fixture.TestSendPeekMessage)
	t.Run("MemoryMessageQueue:Peek No Message", fixture.TestPeekNoMessage)
	t.Run("MemoryMessageQueue:Move To Dead Message", fixture.TestMoveToDeadMessage)
	t.Run("MemoryMessageQueue:On Message", fixture.TestOnMessage)
}
