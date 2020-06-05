import React from 'react';
import { Link } from 'react-router-dom';

function SelectApplication(props) {
  return (
    <>
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
