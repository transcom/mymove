import React from 'react';
import { Link } from 'react-router-dom';
import LoginButton from 'shared/User/LoginButton';

function SelectApplication(props) {
  return (
    <>
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <h1 className="margin-right-2">Select your application</h1>
        <div style={{ listStyle: 'none' }}>
          <LoginButton />
        </div>
      </div>
      <Link to="/moves/queue" onClick={props.handleApplicationSelection}>
        TOO Move Queue
      </Link>
      <br />
      <Link to="/invoicing/queue" onClick={props.handleApplicationSelection}>
        TIO Payment Request Queue
      </Link>
    </>
  );
}

export default SelectApplication;
