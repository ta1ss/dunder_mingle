import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import App from './App';
import ErrorPage from './components/ErrorPage';
import Home from "./components/Home.js";
import Profile from './components/Profile.js';
import Users from './components/Users';
import Groups from './components/Groups';
import Group from './components/Group';
import Login from './components/Login';
import PostDetail from './components/PostDetail';
import Register from './components/Register';
import Welcome from './components/Welcome';
import './style.css';
import Notifications from './components/Notifications';
import Event from './components/Event';

const router = createBrowserRouter([
  { path: '/',
    element: <App />, 
    errorElement: <ErrorPage />, 
    children: [
      {index : true, element: <Home />},
      {path: '/', element: <Home />},
      // {path: '/welcome', element: <Welcome />},
      {path: '/login', element: <Login />},
      {path: '/register', element: <Register />},
      {path: '/profile/me', element: <Profile />},
      {path: '/profile/:id', element: <Profile />},
      {path: '/post/:postId', element: <PostDetail />},
      {path: '/notifications', element: <Notifications />},
      {path: '/users', element: <Users />},
      {path: '/groups', element: <Groups />},
      {path: '/group/:id', element: <Group />},
      {path: '/group/event/:id', element: <Event />},
      {path: '/group/:groupId/post/:postId', element: <PostDetail />},
  ]},
]);


const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
   <RouterProvider router={router}/>
);
