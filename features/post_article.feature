Feature: POST /article

    Scenario: article case1
        Given mysql 执行
            """
            DELETE FROM articles WHERE title IN ('标题1')
            """
        When http 请求 POST /article
            """
            {
                "json": {
                    "title": "标题1",
                    "authorID": 666,
                    "author": "hatlonely",
                    "content": "hello world"
                }
            }
            """
        Then http 检查 201
        Then mysql 检查 "SELECT * FROM articles WHERE title='标题1'"
            """
            {
                "title": "标题1",
                "author_id": 666,
                "author": "hatlonely",
                "content": "hello world"
            }
            """
        Given mysql 执行
            """
            DELETE FROM articles WHERE title IN ('标题1')
            """