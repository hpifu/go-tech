Feature: article 测试

    Scenario: article
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (1,2)
            """
        Given mysql 执行
            """
            DELETE FROM account.accounts WHERE id IN (666)
            """
        Given redis del "godtoken"
        Given mysql 执行
            """
            INSERT INTO articles (id, title, author_id, author, content)
            VALUES (1, "标题1", 666, "hatlonely", "hello world")
            """
        Given mysql 执行
            """
            INSERT INTO articles (id, title, author_id, author, content)
            VALUES (2, "标题2", 666, "hatlonely", "hello world")
            """
        Given mysql 执行
            """
            INSERT INTO account.accounts (id, phone, email, password, first_name, last_name, birthday, gender, avatar)
            VALUES (666, "13112345673", "hatlonely@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1, "hatlonely.png")
            """
        Given redis set string "godtoken"
            """
            1c15b6b0b18aa0d3a5d2de37484f992c
            """
        When http 请求 GET /article/1
        Then http 检查 200
            """
            {
                "json": {
                    "id": 1,
                    "title": "标题1",
                    "authorID": 666,
                    "author": "hatlonely",
                    "content": "hello world"
                }
            }
            """
        When http 请求 GET /article/2
        Then http 检查 200
            """
            {
                "json": {
                    "id": 2,
                    "title": "标题2",
                    "authorID": 666,
                    "author": "hatlonely",
                    "content": "hello world"
                }
            }
            """
        When http 请求 GET /article/3
        Then http 检查 204
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (1,2)
            """
        Given mysql 执行
            """
            DELETE FROM account.accounts WHERE id IN (666)
            """
        Given redis del "godtoken"
