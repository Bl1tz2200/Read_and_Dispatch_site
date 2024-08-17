import { useRef } from 'react'
import './Login.css'
import { API_ADRESS } from '../main'
import { useAuth } from '../Providers/AuthProvider'

export function Login() { // Page where User passes auth and get his JWT token

  var userName = useRef(null)
  var userPassword = useRef(null)
  const { setToken } = useAuth();

  async function logIn() { // Request to backend to get token

    await fetch(
      `${API_ADRESS}/login`,
      {
        method: "Post",
        body: JSON.stringify({
          "UserName": userName.current.value,
          "UserPassword": userPassword.current.value,
        })
      }).then(response => {

        if (response.status == 200) {

          return response.json()

        } else if (response.status == 403) {

          return response.json()

        } else {

          alert("Something went wrong on a server! Try again later")
          window.location.reload()
          return null

        }

      }).then(responseMessage => {
        if (responseMessage) {

          if (responseMessage.token) { // If no errors
            setToken(`${responseMessage.token}`) // Set token to the local storage

            alert("Successfully!")
            window.location.pathname = '/'

          } else { // If error appeared

            if (responseMessage === "User doesn't exist") { // If error appeared because User doesn't exist

              alert("Invalid Username!")
              document.getElementById("usernameIDLogin").className = "username incorrectLogin"

            } else if (responseMessage === "Passwords don't compare") { // If error appeared because User password is different

              alert("Invalid Password!")
              document.getElementById("passwordIDLogin").className = "password incorrectLogin"

            } else { // If error appeared because of something else

              alert("Something went wrong on a server! Try again later")
              window.location.reload()

            }

          }

        }

      })

  }

  return (
    <>
      <div id='rootLogin'>
        <div className="mainPage">
          <h1>Login</h1>
          <div>
            <input ref={userName} className="inputLogin" />
            <i className="username" id="usernameIDLogin">Username</i>
          </div>
          <div>
            <input type="password" ref={userPassword} className="inputPassword" />
            <i className="password" id="passwordIDLogin">Password</i>
          </div>
          <input type="button" className="LogIn" value="Log In" onClick={() => { logIn() }} />
          <h2>Forgot <a href='/reset'><u>Username/Password</u></a>?</h2>
          <h3>Want to <a href='/signup'><u>Sign up</u></a>?</h3>
        </div>
      </div>
    </>
  )
}
