import { createContext, useContext, useEffect, useMemo, useState } from "react";

const AuthContext = createContext(); // Create context

export const AuthProvider = ({ children }) => { // Provider that will set and storage User's JWT token
  const [token, setToken] = useState(localStorage.getItem("token")); // Create token and setToken function, token will have value from local storage

  useEffect(() => { // If token has changed
    if (token) { // If token has changed with some value

      localStorage.setItem('token',token); // Set new token to local storage

    } else { // If token has changed without value

      localStorage.removeItem('token')

    }
  }, [token]);

  const contextValue = useMemo( // Сaching token and setToken function
    () => ({
      token,
      setToken,
    }),
    [token]
  );

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider> // Wrap children DOM with provider that use AuthContext as context
  );
};

export const useAuth = () => { // export provider Context to have access to token and setToken function
  return useContext(AuthContext);
};