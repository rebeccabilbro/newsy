# newsy

`newsy` is a command line news client for Reddit and HackerNews that loads the latest stories, displays them in the terminal, and saves them to a file. `conewsy` is the IO concurrent version.

`searchy` is a simple search engine for news stories.

`findy` is a task parallel search engine to the Reddit and HackerNews client. While new stories are being fetched, users can still be making requests to the server (implemented here via `localhost:8080`).

These program were built as exercises while working through [Hands-on Concurrency with Go by Leo Tindall](https://www.safaribooksonline.com/videos/hands-on-concurrency-with/9781788993746).


## Resources

- [Go wrapper for HackerNews API](https://github.com/caser/gophernews)
- [Create a Reddit app](https://www.reddit.com/prefs/apps)
- [Using Reddit's old cookie authentication method](https://github.com/jzelinskie/geddit#examples)
- [Reddit OAuth method Quickstart](https://github.com/reddit-archive/reddit/wiki/OAuth2-Quick-Start-Example)
- [Python Quickstart](https://github.com/reddit-archive/reddit/wiki/OAuth2-Quick-Start-Example#python-example) ;)
