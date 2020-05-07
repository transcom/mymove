import React from 'react';
import { Link } from 'react-router-dom';

function SelectApplication(props) {
  return (
    <>
      <h1>Select your application</h1>
      <Link to="/moves/queue" onClick={props.handleApplicationSelection}>
        TOO Move Queue
      </Link>
    </>
  );
}

export default SelectApplication;
