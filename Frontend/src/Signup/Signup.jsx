import { useRef } from 'react'
import './Signup.css'
import { API_ADRESS } from '../main'

export function Signup() { // Page where User can sign up on the site
  const userName = useRef(null)
  const password = useRef()
  const passwordRepeat = useRef(null)
  const email = useRef(null)


  async function sendUser() { // Request to backend to registrate User on site

    if (userName.current.value) { // If username has been inputed

      if (password.current.value === passwordRepeat.current.value) { // If passwords are the same

        if (email.current.value.includes("@") && email.current.value.includes(".")) { // If email field includes @ and .

          await fetch(
            `${API_ADRESS}/signup`,
            {
              method: "Post",
              body: JSON.stringify({
                "UserName": userName.current.value,
                "UserPassword": password.current.value,
                "UserEmail": email.current.value
              })
            }).then(response => {

              if (response.status == 201) {

                alert("Successfully!")
                window.location.pathname = '/login'
                return null

              } else if (response.status == 409) {

                return response.json()

              } else {

                alert("Something went wrong on a server! Try again later")
                window.location.reload()
                return null

              }

            }).then(responseMessage => {
              if (responseMessage) { // If user exists

                if (responseMessage === "User already exist with same UserName") { // If exists User with the same username

                  alert("This username are already in use!")
                  document.getElementById("usernameID").className = "usernameSignUp incorrectSignup"

                } else { // If exists User with the same email

                  alert("This email are already in use!")
                  document.getElementById("emailID").className = "email incorrectSignup"

                }
                
              }

            })

        } else { // If email field doesn't have @ and .

          document.getElementById("emailID").className = "email incorrectSignup"

        }

      } else { // If passwords aren't the same
        document.getElementById("passwordRepeatID").className = "repeatpassword incorrectSignup"
      }

    } else { // If username field hasn't been inputed
      document.getElementById("usernameID").className = "usernameSignUp incorrectSignup"
    }

  }


  return (
    <>
      <div id='rootSignUp'>
        <div className="mainPageSignUp">
          <h1>Sign Up</h1>
          <div>
            <input type="text" ref={userName} className="inputLoginSignUp" />
            <i className="usernameSignUp" id="usernameID">Username</i>
          </div>
          <div>
            <input type="password" ref={password} className="inputPasswordSignUp" />
            <i className="passwordSignUp" id="passwordID">Password</i>
          </div>
          <div>
            <input type="password" ref={passwordRepeat} className="inputRepeatPassword" />
            <i className="repeatpassword" id="passwordRepeatID">Repeat Password</i>
          </div>
          <div>
            <input type="email" ref={email} className="inputEmail" />
            <i className="email" id="emailID">Email</i>
          </div>
          <input type="submit" className="SignUp" value="Sign Up" onClick={() => { sendUser() }} />
          <h3>Want to <a href='/login'><u>Log in</u></a>?</h3>
        </div>
      </div>
    </>
  )
}
