package backend

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/datastore"
)

type DelphisBackend interface {
	CreateNewDiscussion(ctx context.Context, creatingUser *model.User, anonymityType model.AnonymityType, title string) (*model.Discussion, error)
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
	GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error)
	GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error)
	SubscribeToDiscussion(ctx context.Context, subscriberUserID string, postChannel chan *model.Post, discussionID string) error
	UnSubscribeFromDiscussion(ctx context.Context, subscriberUserID string, discussionID string) error
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error)
	AssignFlair(ctx context.Context, participant model.Participant, flairID string) (*model.Participant, error)
	CreateFlair(ctx context.Context, userID string, templateID string) (*model.Flair, error)
	GetFlairByID(ctx context.Context, id string) (*model.Flair, error)
	GetFlairsByUserID(ctx context.Context, userID string) ([]*model.Flair, error)
	RemoveFlair(ctx context.Context, flair model.Flair) (*model.Flair, error)
	UnassignFlair(ctx context.Context, participant model.Participant) (*model.Participant, error)
	ListFlairTemplates(ctx context.Context, query *string) ([]*model.FlairTemplate, error)
	CreateFlairTemplate(ctx context.Context, displayName *string, imageURL *string, source string) (*model.FlairTemplate, error)
	RemoveFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error)
	GetFlairTemplateByID(ctx context.Context, id string) (*model.FlairTemplate, error)
	CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Participant, error)
	GetParticipantByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*model.Participant, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	GetParticipantByID(ctx context.Context, id string) (*model.Participant, error)
	CreatePost(ctx context.Context, discussionID string, participantID string, content string) (*model.Post, error)
	NotifySubscribersOfCreatedPost(ctx context.Context, post *model.Post, discussionID string) error
	GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error)
	GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error)
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error)
	CreateUser(ctx context.Context) (*model.User, error)
	GetOrCreateUser(ctx context.Context, input LoginWithTwitterInput) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetViewerByID(ctx context.Context, viewerID string) (*model.Viewer, error)
	GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error)
	CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error)
	GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error)
	UpsertSocialInfo(ctx context.Context, socialInfo model.SocialInfo) (*model.SocialInfo, error)

	NewAccessToken(ctx context.Context, userID string) (*auth.DelphisAccessToken, error)
	ValidateAccessToken(ctx context.Context, token string) (*auth.DelphisAuthedUser, error)
	ValidateRefreshToken(ctx context.Context, token string) (*auth.DelphisRefreshTokenUser, error)
}

type delphisBackend struct {
	db              datastore.Datastore
	auth            auth.DelphisAuth
	cache           cache.ChathamCache
	discussionMutex sync.Mutex
}

func NewDelphisBackend(conf config.Config, awsSession *session.Session) DelphisBackend {
	chathamCache := cache.NewInMemoryCache()
	return &delphisBackend{
		db:              datastore.NewDatastore(conf, awsSession),
		auth:            auth.NewDelphisAuth(&conf.Auth),
		cache:           chathamCache,
		discussionMutex: sync.Mutex{},
	}
}
