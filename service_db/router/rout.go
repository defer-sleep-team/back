package router

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	comm "github.com/defer-sleep-team/Aether_backend/database/internal/comments"
	post "github.com/defer-sleep-team/Aether_backend/database/internal/post"
	subpl "github.com/defer-sleep-team/Aether_backend/database/internal/subscription_plans"
	us "github.com/defer-sleep-team/Aether_backend/database/internal/user"
	usfl "github.com/defer-sleep-team/Aether_backend/database/internal/user_followers"
	usip "github.com/defer-sleep-team/Aether_backend/database/internal/user_ips"
	ussub "github.com/defer-sleep-team/Aether_backend/database/internal/user_subscriptions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Rout(db *sql.DB) error {

	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("never should have pwn here, FSB went 4 u <3")
	})
	api := app.Group("/database_zov_russ_cbo")
	v1 := api.Group("/users")
	v2 := api.Group("/posts")
	v3 := api.Group("/user_ips")
	v4 := api.Group("/subscription_plans")
	v5 := api.Group("/user_followers")
	v6 := api.Group("/subscriptions")
	v7 := api.Group("/comments")
	v1.Post("/", func(c *fiber.Ctx) error {
		user := new(us.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		if err := us.CreateUser(db, user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})
	v1.Get("/get/:iduser/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		idStrUser := c.Params("iduser")
		idUser, err := strconv.Atoi(idStrUser)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		user, err := us.GetUser(db, id, idUser)
		if err != nil {
			log.Print(id)
			if errors.Is(err, us.ErrUserBlocked) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user is blocked"})
			}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		return c.JSON(user)
	})
	v1.Get("/exists/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")
		user, err := us.GetUserByEmail(db, email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		return c.JSON(user)
	})

	v1.Get("/getemail/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		user, err := us.GetEmailUser(db, id)
		if err != nil {
			log.Print(id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})

	v1.Post("/blockuser/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		err = us.BlockUser(db, id)
		if err != nil {
			log.Print(id)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to block user"})
		}

		return c.JSON(fiber.Map{"message": "user blocked successfully"})
	})

	v1.Get("/getblockstatus/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		user, err := us.GetBlockStatus(db, id)
		if err != nil {
			log.Print(id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})
	v1.Get("/getusername/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		user, err := us.GetUsernameUser(db, id)
		if err != nil {
			log.Print(id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})
	v1.Get("/getavatar/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		user, err := us.GetAvatarAndBackgroundUser(db, id)
		if err != nil {
			log.Print(id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})
	v1.Get("/getprivilige/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		user, err := us.GetPrivilegeLevelUser(db, id)
		if err != nil {
			log.Print(id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})
	v1.Put("/changeuser/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, _ := strconv.Atoi(idStr)
		updatedUser := new(us.User)
		if err := c.BodyParser(updatedUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		if err := us.UpdateUser(db, id, updatedUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updatedUser)
	})
	v1.Delete("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, _ := strconv.Atoi(idStr)
		if err := us.DeleteUser(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "user deleted"})
	})
	v1.Post("/auth", func(c *fiber.Ctx) error {
		user := new(us.User)
		newUser := new(us.User)
		var err error
		if err = c.BodyParser(user); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		if newUser, err = us.GetAuthUser(db, *user); err != nil {
			if errors.Is(err, us.ErrUserBlocked) {
				log.Print(err)
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Вы заблокированы на этом сервисе, обратитесь в поддержку"})
			} else if errors.Is(err, sql.ErrNoRows) {
				log.Print(err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неправильный логин или пароль"})
			}
		}
		log.Print(err)
		return c.JSON(newUser)
	})
	// вторая структура в пост go
	v2.Post("/", func(c *fiber.Ctx) error {
		newPost := new(post.Post)
		if err := c.BodyParser(newPost); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		id, err := post.CreatePost(db, *newPost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		newPost.ID = id
		return c.JSON(newPost)
	})
	v2.Post("/full", func(c *fiber.Ctx) error {
		newPostRequest := new(post.IncomingPostRequest)
		if err := c.BodyParser(newPostRequest); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		response, err := post.CreateFullPost(db, *newPostRequest)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(response)
	})
	// Эндпоинт для получения одного поста по ID
	v2.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		postID, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		postDetails, err := post.GetPost(db, postID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(postDetails)
	})

	v2.Get("/ratio/:offset/:n/:id", func(c *fiber.Ctx) error {
		nStr := c.Params("n")
		id := c.Params("id")
		ofst := c.Params("offset")
		postID, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}
		n, err := strconv.Atoi(nStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid number of posts"})
		}
		offset, err := strconv.Atoi(ofst)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid number of posts"})
		}
		posts, err := post.GetNPostsByRatio(db, offset, n, postID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(posts)
	})

	// для получения тегов поста
	v2.Get("/tags/:tagNames", func(c *fiber.Ctx) error {
		tagNamesStr := c.Params("tagNames")
		tagNames := strings.Split(tagNamesStr, ",")

		tagIDs, err := post.TagIDByName(db, tagNames)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(tagIDs)
	})
	//top
	// v2.Get("/topposts", func(c *fiber.Ctx) error {
	// 	postIDs, err := post.GetTopPostsByRatio(db)
	// 	if err != nil {
	// 		log.Print(err)
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get top posts"})
	// 	}

	// 	return c.JSON(fiber.Map{"post_ids": postIDs})
	// })
	// проверка на лайк
	v2.Get("/isliked/:postID/:userID", func(c *fiber.Ctx) error {
		postIDStr := c.Params("postID")
		userIDStr := c.Params("userID")

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}

		isLiked, err := post.IsLiked(db, postID, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"isLiked": isLiked})
	})

	v2.Post("/like/:postID/:userID", func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("postID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		userID, err := strconv.Atoi(c.Params("userID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}

		err = post.LikePost(db, postID, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Post liked successfully"})
	})
	// изменения тегов
	v2.Put("/updateposttags/:postID", func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("postID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		tagNamesStr := c.FormValue("tags")
		if tagNamesStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tags are required"})
		}

		tagNames := strings.Split(tagNamesStr, ",")
		for i, tagName := range tagNames {
			tagNames[i] = strings.TrimSpace(tagName)
		}

		err = post.UpdatePostTags(db, postID, tagNames)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Post tags updated successfully"})
	})

	// Маршрут для удаления лайка с поста
	v2.Delete("/unlike/:postID/:userID", func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("postID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		userID, err := strconv.Atoi(c.Params("userID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}

		err = post.UnlikePost(db, postID, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Post unliked successfully"})
	})

	// comment
	v2.Post("/addcomment/:postID/:userID", func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("postID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		userID, err := strconv.Atoi(c.Params("userID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}

		var request comm.CommentRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
		}

		if request.Content == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "content is required"})
		}

		err = post.AddComment(db, postID, userID, request.Content)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Comment added successfully"})
	})
	// Эндпоинт для получения N постов, на которые подписан пользователь, отсортированных по ratio и новизне
	v2.Get("/subscription/:userID/:n", func(c *fiber.Ctx) error {
		userIDStr := c.Params("userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}

		nStr := c.Params("n")
		n, err := strconv.Atoi(nStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid number of posts"})
		}

		posts, err := post.GetNPostsBySubscription(db, userID, n)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(posts)
	})

	// Эндпоинт для получения 20 рекомендованных постов для пользователя
	v2.Get("/:userid/recommendations/:offset/:n", func(c *fiber.Ctx) error {
		log.Print("/recommendations called")
		userID := c.Params("userid")
		n := c.Params("n")
		ofst := c.Params("offset")
		log.Print(userID, "+", n)
		uid, err := strconv.Atoi(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		num, err := strconv.Atoi(n)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		offset, err := strconv.Atoi(ofst)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		posts, err := post.GetNPostsByRatio(db, offset, num, uid)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(posts)
	})

	v2.Get("/posts_of/:userid/:n", func(c *fiber.Ctx) error {
		log.Print("/posts_of called")
		userID := c.Params("userid")
		n := c.Params("n")
		log.Print(userID, "+", n)
		uid, err := strconv.Atoi(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		num, err := strconv.Atoi(n)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		posts, err := post.GetNPostsOfUser(db, uid, num)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(posts)
	})

	/*v2.Get("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "fail str"})
		}
		post, err := post.GetPost(db, id)
		if err != nil {
			log.Print(id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		return c.JSON(post)
	})*/

	v2.Put("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, _ := strconv.Atoi(idStr)
		updatedPost := new(post.Post)
		if err := c.BodyParser(updatedPost); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		updatedPost.ID = id
		if err := post.UpdatePost(db, *updatedPost); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updatedPost)
	})

	v2.Delete("/:id/:uid", func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}
		userID, err := strconv.Atoi(c.Params("uid"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		err = post.DeletePost(db, postID, userID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Post deleted successfully"})
	})
	v2.Delete("/sudo/delete/post/:id", func(c *fiber.Ctx) error {
		log.Println("222221")
		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		err = post.DeletePostAdmin(db, postID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Post deleted successfully"})
	})
	// изменения поста описания
	v2.Put("/updatepost/:postID", func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("postID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		description := c.FormValue("description")
		if description == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "description is required"})
		}

		err = post.UpdatePostDescription(db, postID, description)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Post description updated successfully"})
	})

	v2.Delete("/deletecomment/:commentID", func(c *fiber.Ctx) error {
		commentID, err := strconv.Atoi(c.Params("commentID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment ID"})
		}

		err = post.DeleteComment(db, commentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Comment deleted successfully"})
	})
	v3.Post("/", func(c *fiber.Ctx) error {
		newUserIP := new(usip.UserIP)
		if err := c.BodyParser(newUserIP); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		err := usip.InsertUserIP(db, *newUserIP)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(newUserIP)
	})

	v3.Get("/:user_id", func(c *fiber.Ctx) error {
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		userIPs, err := usip.GetUserIPs(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(userIPs)
	})

	v3.Put("/:user_id", func(c *fiber.Ctx) error {
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		updatedUserIP := new(usip.UserIP)
		if err := c.BodyParser(updatedUserIP); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		updatedUserIP.UserID = userID
		err = usip.UpdateUserIP(db, *updatedUserIP)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updatedUserIP)
	})

	v3.Delete("/:user_id", func(c *fiber.Ctx) error {
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		err = usip.DeleteUserIP(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "user IPs deleted"})
	})

	v4.Post("/", func(c *fiber.Ctx) error {
		newPlan := new(subpl.SubscriptionPlan)
		if err := c.BodyParser(newPlan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		err := subpl.InsertSubscriptionPlan(db, *newPlan)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(newPlan)
	})

	v4.Get("/:user_id", func(c *fiber.Ctx) error {
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		plans, err := subpl.GetSubscriptionPlans(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(plans)
	})

	v4.Put("/:user_id", func(c *fiber.Ctx) error {
		idStr := c.Params("user_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}
		updatedPlan := new(subpl.SubscriptionPlan)
		if err := c.BodyParser(updatedPlan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		updatedPlan.ID = id
		err = subpl.UpdateSubscriptionPlan(db, *updatedPlan)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updatedPlan)
	})
	// тут на самом деле вопрос
	v4.Delete("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}
		err = subpl.DeleteSubscriptionPlan(db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "subscription plan deleted"})
	})
	// возможно занесение одного и того же человека несколько раз
	v5.Post("/", func(c *fiber.Ctx) error {
		newFollower := new(usfl.UserFollower)
		if err := c.BodyParser(newFollower); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		err := usfl.InsertUserFollower(db, *newFollower)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(newFollower)
	})
	v5.Get("/follower/:followees_id", func(c *fiber.Ctx) error {
		followerIDStr := c.Params("followees_id")
		followerID, err := strconv.Atoi(followerIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid follower ID"})
		}
		followers, err := usfl.GetUserFollowees(db, followerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(followers)
	})
	v5.Get("/:follower_id", func(c *fiber.Ctx) error {
		followerIDStr := c.Params("follower_id")
		followerID, err := strconv.Atoi(followerIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid follower ID"})
		}
		followers, err := usfl.GetUserFollowers(db, followerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(followers)
	})

	v5.Put("/:follower_id/:current_followee_id/:new_followee_id", func(c *fiber.Ctx) error {
		followerIDStr := c.Params("follower_id")
		currentFolloweeIDStr := c.Params("current_followee_id")
		newFolloweeIDStr := c.Params("new_followee_id")

		followerID, err := strconv.Atoi(followerIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid follower ID"})
		}

		currentFolloweeID, err := strconv.Atoi(currentFolloweeIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid current followee ID"})
		}

		newFolloweeID, err := strconv.Atoi(newFolloweeIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid new followee ID"})
		}

		err = usfl.UpdateUserFollower(db, followerID, currentFolloweeID, newFolloweeID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "user follower updated"})
	})

	v5.Delete("/:follower_id/:followee_id", func(c *fiber.Ctx) error {
		followerIDStr := c.Params("follower_id")
		followeeIDStr := c.Params("followee_id")

		followerID, err := strconv.Atoi(followerIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid follower ID"})
		}

		followeeID, err := strconv.Atoi(followeeIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid followee ID"})
		}

		err = usfl.DeleteUserFollower(db, followerID, followeeID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "user follower deleted"})
	})

	v6.Post("/", func(c *fiber.Ctx) error {
		subscription := new(ussub.UserSubscription)
		if err := c.BodyParser(subscription); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		err := ussub.CreateUserSubscription(db, *subscription)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(subscription)
	})

	v6.Get("/:user_id", func(c *fiber.Ctx) error {
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		subscription, err := ussub.GetUserSubscription(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if subscription == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subscription not found"})
		}
		return c.JSON(subscription)
	})

	v6.Put("/:user_id", func(c *fiber.Ctx) error {
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		updatedSubscription := new(ussub.UserSubscription)
		if err := c.BodyParser(updatedSubscription); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		updatedSubscription.UserID = userID
		err = ussub.UpdateUserSubscription(db, *updatedSubscription)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updatedSubscription)
	})

	v6.Delete("/:user_id", func(c *fiber.Ctx) error {
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		err = ussub.DeleteUserSubscription(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "user subscription deleted"})
	})

	v7.Post("/", func(c *fiber.Ctx) error {
		comment := new(comm.Comments)
		if err := c.BodyParser(comment); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		comment.RegDate = time.Now()
		if err := comm.CreateComment(db, comment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(comment)
	})

	v7.Get("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}
		comment, err := comm.GetComment(db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(comment)
	})

	v7.Get("/getallcomments", func(c *fiber.Ctx) error {
		comments, err := comm.GetAllComments(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(comments)
	})

	v7.Put("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}
		updatedComment := new(comm.Comments)
		if err := c.BodyParser(updatedComment); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		updatedComment.RegDate = time.Now()
		if err := comm.UpdateComment(db, id, updatedComment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updatedComment)
	})

	v7.Get("/comments/:postID/:n", func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("postID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
		}

		n, err := strconv.Atoi(c.Params("n"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid number of comments"})
		}

		comments, err := comm.GetCommentsForPost(db, postID, n)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(comments)
	})

	v7.Delete("/:id/:userID", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}
		userIDStr := c.Params("userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}
		if err := comm.DeleteComment(db, id, userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
	v7.Delete("/sudo/delete/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}
		if err := comm.DeleteCommentSudo(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
	err := app.Listen(":8003")
	return err
}
