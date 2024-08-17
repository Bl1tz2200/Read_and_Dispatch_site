import dispatchLogo from '../assets/dispatch_book.svg'
import dispatchAnimation from '../assets/dispatch.gif'
import './User.css'
import { useEffect, useState, useRef, Fragment } from 'react'
import { useAuth } from '../Providers/AuthProvider'
import { API_ADRESS } from "../main.jsx"
import { ShowMiniInfo } from '../PublicFunctions/ShowMiniInfo.jsx'

export function User() { // Page with User info, functions and dispatches

  const { token, setToken } = useAuth();
  const [userData, setUserData] = useState({})
  const [dispatchIds, setdispatchIds] = useState([])

  async function getUserdata(token) { // Request to backend to get data
    const result = await fetch(`${API_ADRESS}/auth`, {
      headers:
      {
        "Authorization": token,
      }
    }).then(response => {

      if (response.status == 202) {

        return response.json()

      } else if (response.status == 401) {

        return response.json()

      } else {

        return null

      }

    }).then(responseMessage => {
      if (responseMessage) {

        if (responseMessage.UserName) { // If no errors

          return { "UserName": responseMessage.UserName, "UserEmail": responseMessage.UserEmail }

        } else { // If auth wasn't passed

          setToken()
          window.location.reload()

        }

      } else { // If error appears

        alert("Something went wrong on a server! Try again later")
        window.location.pathname = "/"

      }

    })

    return await result
  }

  async function changeUsername(token) { // Request to backend to change username
    console.log(typeof (token))

    let input_userName = document.getElementById("userNameUser")

    await fetch(`${API_ADRESS}/changeUN`, {
      method: "POST",
      headers:
      {
        "Authorization": token,
      },
      body: JSON.stringify({
        "UserName": input_userName.value,
      })
    }).then(response => {

      if (response.status == 202) {

        alert("Successfully!")
        setToken()
        window.location.pathname = '/login'

      } else if (response.status == 401) {

        alert("Something went wrong while authorization! Try to log in again")
        setToken()
        window.location.reload()

      } else {

        alert("Something went wrong on a server! Try again later")
        window.location.reload()

      }

    })

  }

  async function changeEmail(token) { // Request to backend to change email
    let input_userEmail = document.getElementById("UserEmail")

    await fetch(`${API_ADRESS}/changeUE`, {
      method: "POST",
      headers:
      {
        "Authorization": token,
      },
      body: JSON.stringify({
        "UserEmail": input_userEmail.value,
      })
    }).then(response => {

      if (response.status == 202) {

        alert("Successfully!")
        window.location.reload()

      } else if (response.status == 401) {

        alert("Something went wrong while authorization! Try to log in again")
        setToken()
        window.location.reload()

      } else {

        alert("Something went wrong on a server! Try again later")
        window.location.reload()

      }

    })

  }

  async function deleteUser(token) { // Request to backend to delete user
    await fetch(`${API_ADRESS}/delUser`, {
      headers:
      {
        "Authorization": token,
      }
    }).then(response => {

      if (response.status == 202) {

        alert("Successfully!")
        setToken()
        window.location.pathname = '/'

      } else if (response.status == 401) {

        alert("Something went wrong while authorization! Try to log in again")
        setToken()
        window.location.reload()

      } else {

        alert("Something went wrong on a server! Try again later")
        window.location.reload()

      }

    })

  }

  async function getUserDispatches(token) { // Request to backend to get user's dispatches
    const ids = await fetch(`${API_ADRESS}/getUserDispatches`, {
      headers:
      {
        "Authorization": token,
      },
    }).then(response => {

      if (response.status == 202) {

        return response.json()

      } else if (response.status == 401) {

        alert("Something went wrong while authorization! Try to log in again")
        setToken()
        window.location.pathname = '/login'

      } else {

        return "error"

      }

    }).then(ids => {
      if (ids) { // If no errors

        setdispatchIds(ids) // Set ids func made by useState

      } else { // If hasn't got ids or appeared error

        if (ids === "error") { // If appeared error
          alert("Something went wrong on a server! Try again later")
          window.location.reload()
        } 

      }
    })

  }

  function logOut() { // Function that removes JWT token from local storage
    setToken()
    window.location.reload()
  }

  useEffect(() => {
    // Functions below will be run on page load
    getUserDispatches(token)

    getUserdata(token).then(
      username => {
        if (username) {
          setUserData(username)
        }
      })
  }, [])

  useEffect(() => {
    // Functions below will be run when userData was changed

    let input_userEmail = document.getElementById("UserEmail")
    let input_userName = document.getElementById("userNameUser")

    if (Object.keys(userData).length != 0) {
      input_userName.value = userData["UserName"]
      input_userEmail.value = userData["UserEmail"]
    } else {
      input_userName.value = "Server error"
      input_userEmail.value = "Server error"
    }

  }, [userData])

  return (
    <>
      <div className="UserField">
        <div>
          <input type="checkbox" className="check" />
          <h1>Your username:</h1>
          <h2>Change username?</h2>
          <input type='text' className="nameInputUser" id="userNameUser" maxLength="20" />
          <input type="button" className="changeUser" value="Change" onClick={() => { changeUsername(token) }} />
        </div>
        <div className='UserEmail'>
          <input type="checkbox" className="check" />
          <h1>Your email:</h1>
          <h2>Change email?</h2>
          <input type='text' className="nameInputUser" id="UserEmail" />
          <input type="button" className="changeUser" value="Change" onClick={() => { changeEmail(token) }} />
        </div>
        <div className='UserOther'>
          <h1>Settings:</h1>
          <input type="button" className="changeUser deleteUser" id="one" value="Delete account" />
          <input type="button" className="changeUser deleteUser" id="two" value="Are you sure?" onClick={() => { deleteUser(token) }} />
          <input type="button" className="changeUser logOutUser" id="three" value="Log Out" />
          <input type="button" className="changeUser logOutUser" id="four" value="Are you sure?" onClick={() => { logOut() }} />
          <a href='/reset'><input type="button" className="changeUser changePasswordOther" value="Change password" /></a>
        </div>
      </div>
      <div className="headerUser">
        <a href="/dispatch"><img src={dispatchLogo} className="logoUser" alt="Dispatch"
          onMouseOver={mySrc => { mySrc.currentTarget.src = dispatchAnimation; mySrc.currentTarget.className = "animatedLogoUser" }}
          onMouseOut={mySrc => { mySrc.currentTarget.src = dispatchLogo; mySrc.currentTarget.className = "logoUser" }} /></a>
        <input type='text' placeholder='Search in your works...' id="searchUsers" onKeyDown={(enteredKey) => {
          if (enteredKey.key == "Enter"){

            var enteredText = document.getElementById("searchUsers").value

            dispatchIds.map(id =>
            {
  
              if(!(`${document.getElementById(`${id}`).textContent}`.includes(enteredText))){
                document.getElementById(`${id}/div`).style.display = "none"
              } else {
                document.getElementById(`${id}/div`).style.display = "inline-flex"
              }

            }
            )
          }
        }} />
        <div className="userNameUser">
          <h1><a href="/">Go to <br />main page</a></h1>
        </div>
      </div>
      <div className="articlesUser" >
        <ul>
          { // Display all mini dispatches
            dispatchIds.map(id =>
              <Fragment key={id}>
                <a href={`/dispatch/${id}`}>< ShowMiniInfo id={id} /></a>
              </Fragment>
            )
          }
        </ul>
      </div>
    </>
  )
}
