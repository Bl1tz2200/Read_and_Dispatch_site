import React from 'react'
import { Dispatcher } from '../Dispatcher/Dispatcher.jsx'
import { Dispatch } from '../Dispatch/Dispatch.jsx'
import { Login } from '../Login/Login.jsx'
import { Signup } from '../Signup/Signup.jsx'
import { Reset } from '../Reset/Reset.jsx'
import { Changer } from '../Changer/Changer.jsx'
import { User } from '../UserPage/User.jsx'
import { createBrowserRouter, Navigate, Outlet, RouterProvider } from "react-router-dom";
import { useAuth } from '../Providers/AuthProvider.jsx'



const ProtectedRoute = () => { // Function returns page were will go user if he tries to open page that requires auth
    const { token } = useAuth();

    if (!token) {
        return <Navigate to="/login" />;
    }

    return <Outlet />;
};

export const Routes = () => { // All routes
    const { token } = useAuth();

    const routesForPublic = [ // Public pages
        {
            path: "/",
            element: <Dispatcher />,
        },
        {
            path: "/dispatch/:id",
            element: <Dispatch />,
        },
        {
            path: "/reset",
            element: <Reset />,
        },
        {
            path: "/password-change",
            element: <Changer />,
            errorElement: <Dispatcher />,
        },
    ];

    const routesForAuthenticatedOnly = [ // Pages that requires auth
        {
            path: "/",
            element: <ProtectedRoute />,
            children: [
                {
                    path: "/",
                    element: <Dispatcher />,
                },
                {
                    path: "/user",
                    element: <User />,
                },
                {
                    path: "/dispatch/:id",
                    element: <Dispatch />,
                },
                {
                    path: "/dispatch",
                    element: <Dispatch />,
                },
            ],
        },
    ];

    const routesForNotAuthenticatedOnly = [ // Non-Auth Pages
        {
            path: "/login",
            element: <Login />,
        },
        {
            path: "/signup",
            element: <Signup />,
        },
    ];

    const router = createBrowserRouter([ // Creating router
        ...routesForPublic,
        ...(!token ? routesForNotAuthenticatedOnly : []), // If you authed you won't see Non-Auth Pages
        ...routesForAuthenticatedOnly,
    ]);

    return <RouterProvider router={router} />;

};