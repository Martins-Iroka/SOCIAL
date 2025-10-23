package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/Martins-Iroka/social/internal/store"
)

var usernames = []string{
	"CodeMaster7", "GoGopherFan", "TechWiz_22", "AquaDev", "StarCoderX", "PixelPanda", "TheDataSage", "BitStreamer", "LogicLover", "SyntaxSlayer",
	"BinaryBard", "NetNomad", "KeyBoardKatt", "AlgorithmicAce", "FunctionFiend", "Compile_King", "VectorVixen", "ParallelPete", "SwiftStruct", "GoroutineGuru",
	"LambdaLlama", "MacroMaestro", "JitterBugJr", "OpenSourceOllie", "CloudCrusader", "SystemSorcerer", "DynamicDiva", "ImmutableIdiot", "PointerPal", "StackOverFlow",
	"ThreadThriller", "WebWeaver_88", "KernelKnight", "ShellShocked", "GizmoGuru", "CyberCoyote", "DigitalDuke", "ElectricEcho", "FirewallFox", "GridGuardian",
	"HackerHedgehog", "InfraInventor", "JungleJumper", "KiloKiller", "LoopLeaper", "MatrixMan", "NodeNinja", "OracleOwl", "ProtocolPro", "QueryQueen",
	"RouterRabbit", "ScriptSavior", "TunnelTiger", "UpgradeUnicorn", "VirtualVoyager", "WidgetWarrior", "XenonXpert", "YottaYielder", "ZetaZenith", "GigaGamerGirl",
	"NanoNerdBoy", "TeraTypist", "PetaProgrammer", "ExaExplorer", "ZettaZoomer", "YottaYacht", "RustRider", "PythonPirate", "JavaJockey", "CSScultist",
	"HTMLhero", "JScriptJester", "RubyRogue", "SwiftSorcerer", "KotlinKing", "CPlusChief", "ScalaScout", "HaskellHopper", "LispLegend", "PerlPioneer",
	"LuaLover", "DartDynamo", "RacerGo", "CompileTime", "RunTimeRiot", "DebugDancer", "RefactorReady", "TestDrivenTed", "MuxMaster", "ServerSideSam",
	"ClientCalm", "APIAdventurer", "DataLakeDiver", "StreamSailor", "BatchBouncer", "CacheCommando", "DeployDestroyer", "MonitorMage", "SentrySeven", "GuardianGo",
}

var contents = []string{
	// Row 1
	"Just shipped the latest feature! 🚀 #GoLang #SoftwareDev",
	"Loving the new VS Code theme. Productivity just went up 📈.",
	"Why is it always the semicolon? 🤔 Debugging life.",
	"Coffee + Code = A perfect Saturday morning. ☕",
	"Read 'The Go Programming Language' this week. Highly recommend!",
	"Learning about Goroutines today. Concurrency is key! ✨",
	"Setup complete. Time to build something amazing from scratch.",
	"Anyone else think `fmt.Println` is the ultimate debugger?",
	"Just committed my first pull request to an open-source project!",
	"Struggling with channel deadlock... Send help! 🆘",

	// Row 2
	"Design patterns review: The Factory Method is my favorite.",
	"Celebrating 1 year of coding consistently! 🎉 Keep building.",
	"Found a great tutorial on Kubernetes networking. Link in bio!",
	"Is Vim better than Emacs? Discuss. ⚔️",
	"My desk setup got a major upgrade. Feeling inspired! ",
	"Thinking about transitioning to full-stack development. Any advice?",
	"Success! Fixed the race condition that's been bugging me for days.",
	"The power of clean code is truly unmatched. Refactor time.",
	"Exploring machine learning with Go's Gonum library. Fascinating!",
	"Remember to take breaks! Mental health > Merge conflicts. 🧘",
}

var titles = []string{
	"A Deep Dive into Go's Concurrency Model",
	"Optimizing Database Queries in a Microservices Architecture",
	"From Zero to Deploy: Building a REST API in Go",
	"The Hidden Costs of Premature Optimization",
	"10 Essential Tools for Every Gopher's Toolkit",
	"Mastering Custom Error Handling in Go",
	"Deciphering Kubernetes: A Beginner's Guide",
	"Why I Switched from Python to Go for Backend Development",
	"Refactoring Legacy Code: A Step-by-Step Approach",
	"Securing Your Application with JWT and Go",
	"Understanding the Go Garbage Collector",
	"Implementing the Repository Pattern with GORM",
	"Scaling Your Go Application with Goroutines and Channels",
	"My Favorite VS Code Extensions for Go Development",
	"The Power of Interfaces in Clean Go Design",
	"Performance Benchmarking Your Code in Go",
	"Exploring Dependency Injection in Modern Go",
	"The State of WebAssembly in 2024",
	"Writing Testable Code: Beyond Unit Tests",
	"Automating Infrastructure with Terraform and Go",
}

var tags = []string{
	"#golang", "#programming", "#softwaredevelopment", "#tech", "#codinglife", "#devops", "#cloud", "#kubernetes", "#microservices", "#opensource",
	"#webdev", "#frontend", "#backend", "#fullstack", "#datascientist", "#ai", "#machinelearning", "#javascript", "#python", "#rustlang",
}

var commentSlice = []string{
	"Great explanation! This really clarified the concept of closures for me. 👍",
	"I ran into the same bug last week! Did you try rolling back the library version?",
	"Awesome post. Can you share the GitHub repo for the code samples?",
	"Totally agree about the semicolon struggle! 😂 We've all been there.",
	"This is a fantastic use case for Go's error handling patterns. Well done!",
	"What kind of processor are you running? That compile time is insane!",
	"Highly underrated feature of that framework. Thanks for pointing it out.",
	"A solid argument for using interfaces! Clean code FTW.",
	"First!",
	"I think the time complexity of that algorithm could be reduced with a hash map.",
	"Nice color scheme! Is that the Monokai theme with some tweaks?",
	"Could you elaborate on the security implications of this approach?",
	"Inspired! Going to try implementing this in my next personal project.",
	"Vim is the superior editor, change my mind. 😉",
	"The data visualization in your last slide was perfect. Very clear.",
	"This is going straight into my 'must-read' bookmarks. Thanks for sharing!",
	"That feeling when the tests pass on the first try! Pure bliss. ✨",
	"Any recommendations for a good book on advanced concurrency in Go?",
	"My team just adopted this exact pattern. It's saved us a ton of headaches.",
	"Thanks for the reminder to take a break! Going for a quick walk now. 🚶",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.User.CreateUser(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Post.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comment.CreateComment(ctx, comment); err != nil {
			log.Println("Error creating comments:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range num {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "12345",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := range num {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := range num {
		userID := users[rand.Intn(len(users))].ID
		postID := posts[rand.Intn(len(posts))].ID

		comments[i] = &store.Comment{
			PostID:  postID,
			UserID:  userID,
			Content: commentSlice[rand.Intn(len(commentSlice))],
		}
	}

	return comments
}
