import React from 'react';
import { connect } from 'react-redux';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import RequestAccountForm from 'components/Office/RequestAccountForm/RequestAccountForm';

export const RequestAccount = () => {
  const initialValues = {};

  const handleCancel = () => {};

  const handleSubmit = async () => {};

  return (
    <GridContainer>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <RequestAccountForm initialValues={initialValues} onCancel={handleCancel} onSubmit={handleSubmit} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

RequestAccount.propTypes = {};

const mapDispatchToProps = {};

const mapStateToProps = () => ({});

export default connect(mapStateToProps, mapDispatchToProps)(RequestAccount);
