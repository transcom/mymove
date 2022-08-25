import React from 'react';
import { useHistory, useParams } from 'react-router';
import { Button, Grid, GridContainer } from '@trussworks/react-uswds';

const Forbidden = () => {
  const history = useHistory();
  const { moveCode } = useParams();
  // TODO need to figure out how to style this thing
  // TODO grid kinda works buuut it has some awkward breakpoints and really
  // TODO i just need to make a dang box a fixed width and center it, seems like
  // TODO a lot of extra markup for that.
  return (
    <GridContainer>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <h1>Sorry, you can&apos;t access this page</h1>
          <p>This page is only available to authorized users</p>
          <p>
            You are not signed in to MilMove in a role that gives you access. If you believe you should have access,
            contact your administrator.
          </p>
          <Button onClick={() => history.push(`/moves/${moveCode}/details`)}>Go to move details</Button>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export default Forbidden;
