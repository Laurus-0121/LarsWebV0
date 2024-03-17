function submitForm() {
    var username = document.getElementById('user_name').value;
    var password = document.getElementById('pass_word').value;

    var user = {
        user_name: username,
        pass_word: password
    };

    var jsonUser = JSON.stringify(user);
    console.log(jsonUser);  // 在控制台打印 JSON 格式的用户数据，方便调试

    // TODO: 可以将 jsonUser 发送到服务器端

    // TODO: 可以将 jsonUser 发送到服务器端
    fetch('http://localhost:8080/user/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: jsonUser
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log(data); // 服务器返回的数据
            // 登录成功后重定向到 user.html 页面
            localStorage.setItem('Authorization', data.data); // 将 JWT 存储在 localStorage 中
            localStorage.setItem('user',user)
            window.location.href = 'userInfo';

        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });

}

function register() {
    var username = document.getElementById('register_name').value;
    var password = document.getElementById('register_password').value;

    var user = {
        username: username,
        password: password
    };

    var jsonUser = JSON.stringify(user);
    console.log(jsonUser);  // 在控制台打印 JSON 格式的用户数据，方便调试

    // TODO: 可以将 jsonUser 发送到服务器端
    fetch('http://localhost:8080/user/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: jsonUser
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log(data); // 服务器返回的数据
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
}