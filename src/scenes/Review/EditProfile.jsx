import React, { Component, Fragment } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { push } from 'react-router-redux';
import { Field, reduxForm } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';
import SaveCancelButtons from './SaveCancelButtons';
import { updateServiceMember } from 'scenes/ServiceMembers/ducks';
import { moveIsApproved, isPpm } from 'scenes/Moves/ducks';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import { editBegin, editSuccessful, entitlementChangeBegin, entitlementChanged, checkEntitlement } from './ducks';
import scrollToTop from 'shared/scrollToTop';

import './Review.css';
import profileImage from './images/profile.png';

const editProfileFormName = 'edit_profile';

let EditProfileForm = props => {
  const {
    schema,
    handleSubmit,
    submitting,
    valid,
    moveIsApproved,
    initialValues,
    schemaAffiliation,
    schemaRank,
    serviceMember,
  } = props;
  const currentStation = get(serviceMember, 'current_station');
  const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');
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
      {!moveIsApproved && (
        <Fragment>
          <SwaggerField fieldName="affiliation" swagger={schema} required />
          <SwaggerField fieldName="rank" swagger={schema} required />
          <SwaggerField fieldName="edipi" swagger={schema} required />
          <Field name="current_station" component={DutyStationSearchBox} />
        </Fragment>
      )}
      {moveIsApproved && (
        <Fragment>
          <div>
            To change the fields below, contact your local PPPO office at {get(currentStation, 'name')}{' '}
            {stationPhone ? ` at ${stationPhone}` : ''}.
          </div>
          <label>Branch</label>
          <strong>{schemaAffiliation['x-display-value'][initialValues.affiliation]}</strong>
          <label>Rank</label>
          <strong>{schemaRank['x-display-value'][initialValues.rank]}</strong>
          <label>DoD ID #</label>
          <strong>{initialValues.edipi}</strong>

          <label>Current Duty Station</label>
          <strong>{get(initialValues, 'current_station.name')}</strong>
        </Fragment>
      )}
      <SaveCancelButtons valid={valid} submitting={submitting} />
    </form>
  );
};
const validateProfileForm = validateAdditionalFields(['current_station']);
EditProfileForm = reduxForm({
  form: editProfileFormName,
  validate: validateProfileForm,
})(EditProfileForm);

class EditProfile extends Component {
  updateProfile = fieldValues => {
    fieldValues.current_station_id = fieldValues.current_station.id;
    if (fieldValues.rank !== this.props.serviceMember.rank) {
      this.props.entitlementChanged();
    }
    const moveId = this.props.move.id;
    return this.props.updateServiceMember(fieldValues).then(() => {
      // This promise resolves regardless of error.
      if (!this.props.hasSubmitError) {
        this.props.editSuccessful();
        this.props.history.goBack();
        if (this.props.isPpm) {
          this.props.checkEntitlement(moveId);
        }
      } else {
        scrollToTop();
      }
    });
  };

  componentDidMount() {
    this.props.editBegin();
    this.props.entitlementChangeBegin();
  }

  render() {
    const { error, schema, serviceMember, moveIsApproved, schemaAffiliation, schemaRank } = this.props;

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
            moveIsApproved={moveIsApproved}
            schemaRank={schemaRank}
            schemaAffiliation={schemaAffiliation}
            serviceMember={serviceMember}
          />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    serviceMember: get(state, 'serviceMember.currentServiceMember'),
    move: get(state, 'moves.currentMove'),
    error: get(state, 'serviceMember.error'),
    hasSubmitError: get(state, 'serviceMember.hasSubmitError'),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    moveIsApproved: moveIsApproved(state),
    isPpm: isPpm(state),
    schemaRank: get(state, 'swaggerInternal.spec.definitions.ServiceMemberRank', {}),
    schemaAffiliation: get(state, 'swaggerInternal.spec.definitions.Affiliation', {}),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      push,
      updateServiceMember,
      editBegin,
      entitlementChangeBegin,
      editSuccessful,
      entitlementChanged,
      checkEntitlement,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(EditProfile);
