import logo from './logo.svg';
import './App.css';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Card, Divider, Input, message, Space } from 'antd';

type DoneReply = {
  ok: boolean;
  message: string;
};

type AppProps = {
  backendURL: string;
};

function App(props: AppProps) {
  const navigate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();


  const isLoggedIn = !!localStorage.getItem("uc");
  if (isLoggedIn) {
    // TODO: redirect to dashboard, pass in the backendURL as well
  }
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const onUsernameChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
    setUsername(e.target.value);
  }
  const onPasswordChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
    setPassword(e.target.value);
  }

  const onLogin = async () => {
    // TODO: add try catch here
    const res = await fetch(`http://${props.backendURL}/admin/login`, { method: 'POST', mode: 'cors', body: JSON.stringify({ username, password }) });
    let body = await res.json();
    if (body.ok) {
      messageApi.open({
        content: body.message,
        type: 'success',
      })
      navigate("/dashboard", { replace: true });
    } else {
      messageApi.open({
        content: body.message,
        type: 'error',
      })
    }
  }

  return (
    <div className='App' >
      <Card title='Chrononut'>
          <Space>
            <Input value={username} onChange={onUsernameChange} id='uname' placeholder='Username' />
            <Input value={password} onChange={onPasswordChange} id='pwd' type='password' placeholder='Password' />
          </Space>
        <Divider />
            <Button onClick={onLogin}>
              Log In
            </Button>
      </Card>
    </div>
  );
}

export default App;
