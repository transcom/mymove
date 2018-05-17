import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { push } from 'react-router-redux';
import { Field, reduxForm } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';

import { updateServiceMember } from 'scenes/ServiceMembers/ducks';

import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import './Review.css';
import profileImage from './images/profile.png';

const editProfileFormName = 'edit_profile';

let EditProfileForm = props => {
  const { onCancel, schema, handleSubmit, pristine, submitting, valid } = props;
  return (
    <form onSubmit={handleSubmit}>
      <img src={profileImage} alt="" /> Profile
      <hr />
      <h3 className="sm-heading">Edit Profile:</h3>
      <SwaggerField fieldName="first_name" swagger={schema} required />
      <SwaggerField fieldName="middle_name" swagger={schema} />
      <SwaggerField fieldName="last_name" swagger={schema} required />
      <SwaggerField fieldName="suffix" swagger={schema} />
      <hr className="spacer" />
      <SwaggerField fieldName="affiliation" swagger={schema} required />
      <SwaggerField fieldName="rank" swagger={schema} required />
      <SwaggerField fieldName="edipi" swagger={schema} required />
      <Field name="current_station" component={DutyStationSearchBox} />
      <button type="submit" disabled={pristine || submitting || !valid}>
        Save
      </button>
      <button type="button" disabled={submitting} onClick={onCancel}>
        Cancel
      </button>
    </form>
  );
};
const validateProfileForm = validateAdditionalFields(['current_station']);
EditProfileForm = reduxForm({
  form: editProfileFormName,
  validate: validateProfileForm,
})(EditProfileForm);

class EditProfile extends Component {
  returnToReview = () => {
    const reviewAddress = `/moves/${this.props.match.params.moveId}/review`;
    this.props.push(reviewAddress);
  };

  updateProfile = (fieldValues, something, elses) => {
    fieldValues.current_station_id = fieldValues.current_station.id;

    this.props.updateServiceMember(fieldValues).then(() => {
      // This promise resolves regardless of error.
      if (!this.props.error) {
        this.returnToReview();
      }
    });
  };

  render() {
    const { error, schema, serviceMember } = this.props;

    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditProfileForm
            initialValues={serviceMember}
            onSubmit={this.updateProfile}
            onCancel={this.returnToReview}
            schema={schema}
          />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    serviceMember: get(state, 'loggedInUser.loggedInUser.service_member'),
    move: get(state, 'moves.currentMove'),
    error: get(state, 'serviceMember.error'),
    schema: get(
      state,
      'swagger.spec.definitions.CreateServiceMemberPayload',
      {},
    ),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, updateServiceMember }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(EditProfile);
