package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
	allParticipants, err := d.GetParticipantsByDiscussionID(ctx, discussionID)

	if err != nil {
		return nil, err
	}

	viewerObj, err := d.CreateViewerForDiscussion(ctx, discussionID, userID)

	if err != nil {
		return nil, err
	}

	participantObj := model.Participant{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ParticipantID: len(allParticipants),
		DiscussionID:  &discussionID,
		ViewerID: &viewerObj.ID,
		UserID:   &userID,
	}

	_, err = d.db.PutParticipant(ctx, participantObj)

	if err != nil {
		return nil, err
	}

	return &participantObj, nil
}

func (d *delphisBackend) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	return d.db.GetParticipantsByDiscussionID(ctx, id)
}

func (d *delphisBackend) GetParticipantByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
	return d.db.GetParticipantByDiscussionIDUserID(ctx, discussionID, userID)
}

func (d *delphisBackend) GetParticipantByID(ctx context.Context, id string) (*model.Participant, error) {
	participant, err := d.db.GetParticipantByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return participant, nil
}

func (d *delphisBackend) AssignFlair(ctx context.Context, participant model.Participant, flairID string) (*model.Participant, error) {
	return d.db.AssignFlair(ctx, participant, &flairID)
}

func (d *delphisBackend) UnassignFlair(ctx context.Context, participant model.Participant) (*model.Participant, error) {
	return d.db.AssignFlair(ctx, participant, nil)
}

