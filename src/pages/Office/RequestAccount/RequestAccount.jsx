import React from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import RequestAccountForm from 'components/Office/RequestAccountForm/RequestAccountForm';

export const RequestAccount = () => {
  const navigate = useNavigate();

  const initialValues = {};

  const handleCancel = () => {
    navigate(-1);
  };

  const handleSubmit = async () => {
    navigate(-1);
  };

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
