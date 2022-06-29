import React from "react";
import { Form, Input, Button, message } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import { Link } from "react-router-dom";
import axios from "axios";

import { BASE_URL } from "../constants";

function Login(props) {
    const { handleLoggedIn } = props;

    const onFinish = (values) => {
        console.log('Received values of form: ', values);
        //step 1: get username/password from the form
        //step 2: send user info to the server
        //step 3: parse the response
        //  case 1: success -> pass token to App
        //  case 2: fail -> inform user
        const { username, password } = values;
        const opt = {
            method: "POST",
            url: `${BASE_URL}/signin`,
            data: {
                username: username,
                password: password
            },
            headers: { "Content-Type": "application/json" }
        };
        axios(opt)
            .then((res) => {  //case 1: success
                if (res.status === 200) {
                    const { data } = res;
                    //send data to App
                    handleLoggedIn(data);
                    message.success("Login succeed! ");
                }
            })
            .catch((err) => {       //case 2: fail
                console.log("login failed: ", err.message);
                message.error("Login failed!");
            });
    };

    return (
        <Form
            name="normal_login"
            className="login-form"
            onFinish={onFinish}
        >
            <Form.Item
                name="username"
                rules={[
                    {
                        required: true,
                        message: 'Please input your Username',
                    },
                ]}
            >
                <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Username" />
            </Form.Item>
            <Form.Item
                name="password"
                rules={[
                    {
                        required: true,
                        message: 'Please input your Password',
                    },
                ]}
            >
                <Input
                    prefix={<LockOutlined className="site-form-item-icon" />}
                    type="password"
                    placeholder="Password"
                />
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit" className="login-form-button">
                    Log in
                </Button>
                Or <a href="">register now!</a>
            </Form.Item>
        </Form>
    );
}

export default Login;