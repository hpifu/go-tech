Feature: /like/:id

    Scenario: like
        Given mysql 执行 "DELETE FROM likeviews WHERE id IN (1)"
        When http 请求 POST /like/1
        Then http 检查 200
        When http 请求 POST /like/1
        Then http 检查 200
        Then mysql 检查 "SELECT * FROM likeviews WHERE id=1"
            """
            {
                "like": 2,
                "view": 0
            }
            """
        Given mysql 执行 "DELETE FROM likeviews WHERE id IN (1)"
