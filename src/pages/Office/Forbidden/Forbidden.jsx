import React from 'react';
import { useHistory, useParams } from 'react-router';
import { Button } from '@trussworks/react-uswds';

const Forbidden = () => {
  const history = useHistory();
  const { moveCode } = useParams();
  return (
    <div>
      <h1>Sorry, you can&apos;t access this page</h1>
      <p>This page is only available to authorized users</p>
      <p>
        You are not signed in to MilMove in a role that gives you access. If you believe you should have access, contact
        your administrator.
      </p>
      <Button onClick={() => history.push(`/moves/${moveCode}/details`)}>Go to move details</Button>
    </div>
  );
};

export default Forbidden;
