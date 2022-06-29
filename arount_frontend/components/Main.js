import React from 'react';
import { Route, Switch, Redirect } from 'react-router-dom';

import Login from './Login';
import Register from "./Register";
import Home from "./Home";

function Main(props) {

    const { isLoggedIn, handleLoggedIn } = props;

    const showLogin = () => {
        //case 1: isLoggedIn => show home
        //case 2: !isLoggedIn => show Login
        return isLoggedIn
            ?
            <Redirect to="/home" />
            :
            <Login handleLoggedIn={handleLoggedIn}/>
    }

    const showHome = () => {
        //case 1: isLoggedIn => show home
        //case 2: !isLoggedIn => show Login
        return isLoggedIn ? <Home /> : <Redirect to="/login" />;
    };

    return (
        <div className= "main">
            <Switch>
                <Route path="/" exact render={showLogin} />
                <Route path="/login" render={showLogin} />
                <Route path="/register" component={Register} />
                <Route path="/home" render={showHome} />
            </Switch>
        </div>
    );
}

export default Main;