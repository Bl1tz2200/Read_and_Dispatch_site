import { API_ADRESS } from '../main';
import './Reset.css'
import { useNavigate } from "react-router-dom";


export function Reset() { // Page where you should enter email to get password-reset letter
  const navigate = useNavigate();

  async function sendPasswordRecovery() { // request for letter to backend
    var userEmail = document.getElementById("inputEmail")

    if (userEmail.value.includes("@") && userEmail.value.includes(".")) { // If entered email has @ and .

      await fetch(`${API_ADRESS}/sendPasswordRecovery`, {
        method: "POST",
        body: JSON.stringify({
          "UserEmail": userEmail.value,
        })
      }).then(response => {

        if (response.status == 202) {

          alert("Successfully!")
          window.location.pathname = "/login"

        } else {

          alert("Something went wrong on a server! Try again later")
          window.location.reload()

        }

      })
    } else { // If entered email doesn't looks like email
      document.getElementById("inputEmailText").className = "emailReset incorrectReset"
    }

  }

  return (
    <>
      <div id='rootReset'>
        <div className="mainPageReset">
          <h1>Change password</h1>
          <div>
            <input type="email" className="inputEmailReset" id="inputEmail" />
            <i className="emailReset" id="inputEmailText">Email</i>
          </div>
          <input type="button" className="Reset" value="Change" onClick={() => { sendPasswordRecovery() }}/>
          <h4><a onClick={() => navigate(-1)}><u>Go back</u></a></h4>
        </div>
      </div>
    </>
  )
}
