package pubsub

import (
	"github.com/segmentio/kafka-go"
	"sort"
)

type Balancer struct{}

func (r Balancer) ProtocolName() string {
	return "roundrobin"
}

func (r Balancer) UserData() ([]byte, error) {
	return nil, nil
}

func (r Balancer) AssignGroups(members []kafka.GroupMember, topicPartitions []kafka.Partition) kafka.GroupMemberAssignments {
	assignor := newGroupMemberAssignor(members)
	for _, partition := range topicPartitions {
		assignor.Assign(partition)
	}
	groupAssignments := assignor.GetGroupAssignments()

	return groupAssignments
}

type groupMemberAssignor struct {
	memberHolders []*groupMemberHolder
}

func newGroupMemberAssignor(members []kafka.GroupMember) *groupMemberAssignor {
	memberHolders := []*groupMemberHolder{}
	for _, member := range members {
		memberHolders = append(memberHolders, newGroupMemberHolder(member))
	}

	return &groupMemberAssignor{memberHolders: memberHolders}
}

func (a *groupMemberAssignor) Assign(partition kafka.Partition) {
	sort.Sort(byHolder(a.memberHolders))

	for _, holder := range a.memberHolders {
		if holder.hasTopic(partition.Topic) {
			holder.assign(partition)
			return
		}
	}
}

func (a *groupMemberAssignor) GetGroupAssignments() kafka.GroupMemberAssignments {
	// GroupMemberAssignments holds MemberID => topic => partitions

	groupAssignments := kafka.GroupMemberAssignments{}
	for _, holder := range a.memberHolders {
		groupAssignments[holder.getMemberID()] = holder.getAssignments()
	}
	return groupAssignments
}

type groupMemberHolder struct {
	member     *kafka.GroupMember
	partitions []*kafka.Partition
}

func newGroupMemberHolder(member kafka.GroupMember) *groupMemberHolder {
	return &groupMemberHolder{
		member:     &member,
		partitions: []*kafka.Partition{},
	}
}

func (h *groupMemberHolder) hasTopic(targetTopic string) bool {
	for _, topic := range h.member.Topics {
		if topic == targetTopic {
			return true
		}
	}
	return false
}

func (h *groupMemberHolder) assign(partition kafka.Partition) {
	h.partitions = append(h.partitions, &partition)
}

func (h *groupMemberHolder) getMemberID() string {
	return h.member.ID
}

func (h *groupMemberHolder) getAssignments() map[string][]int {
	// topic => partitions
	assignments := map[string][]int{}
	for _, partition := range h.partitions {
		partitions, ok := assignments[partition.Topic]
		if !ok {
			partitions = []int{}
		}
		partitions = append(partitions, partition.ID)
		assignments[partition.Topic] = partitions
	}
	return assignments
}

type byHolder []*groupMemberHolder

func (h byHolder) Len() int {
	return len(h)
}

func (h byHolder) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h byHolder) Less(i, j int) bool {
	return len(h[i].partitions) < len(h[j].partitions)
}

func (h *groupMemberHolder) Len() int {
	return len(h.partitions)
}
