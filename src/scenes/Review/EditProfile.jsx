import React, { Component, Fragment } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { push } from 'connected-react-router';
import { Field, reduxForm } from 'redux-form';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';
import SaveCancelButtons from './SaveCancelButtons';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import { editBegin, editSuccessful, entitlementChangeBegin, entitlementChanged, checkEntitlement } from './ducks';
import scrollToTop from 'shared/scrollToTop';
import {
  selectServiceMemberFromLoggedInUser,
  selectMoveIsSubmitted,
  selectCurrentMove,
  selectHasCurrentPPM,
} from 'store/entities/selectors';

import './Review.css';
import profileImage from './images/profile.png';
import SectionWrapper from 'components/Customer/SectionWrapper';
import ServiceInfoTable from 'components/Customer/Review/ServiceInfoTable';

const editProfileFormName = 'edit_profile';

let EditProfileForm = (props) => {
  const { schema, handleSubmit, submitting, valid, moveIsSubmitted, initialValues, serviceMember } = props;
  const currentStation = get(serviceMember, 'current_station');
  const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');
  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <img src={profileImage} alt="" />{' '}
            <h1
              style={{
                display: 'inline-block',
                marginLeft: 10,
                marginBottom: 16,
                marginTop: 20,
              }}
            >
              Profile
            </h1>
            <SectionWrapper>
              <h2>Edit Profile:</h2>
              <SwaggerField fieldName="first_name" swagger={schema} required />
              <SwaggerField fieldName="middle_name" swagger={schema} />
              <SwaggerField fieldName="last_name" swagger={schema} required />
              <SwaggerField fieldName="suffix" swagger={schema} />
              <hr className="spacer" />
              {!moveIsSubmitted && (
                <>
                  <SwaggerField fieldName="affiliation" swagger={schema} required />
                  <SwaggerField fieldName="rank" swagger={schema} required />
                  <SwaggerField fieldName="edipi" swagger={schema} required />
                  <Field name="current_station" title="Current duty station" component={DutyStationSearchBox} />
                </>
              )}
              {moveIsSubmitted && (
                <ServiceInfoTable
                  firstName={initialValues.first_name}
                  lastName={initialValues.last_name}
                  currentDutyStationName={currentStation.name}
                  currentDutyStationPhone={stationPhone}
                  affiliation={initialValues.affiliation}
                  rank={initialValues.rank}
                  edipi={initialValues.edipi}
                />
              )}
            </SectionWrapper>
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};
const validateProfileForm = validateAdditionalFields(['current_station']);
EditProfileForm = reduxForm({
  form: editProfileFormName,
  validate: validateProfileForm,
})(EditProfileForm);

class EditProfile extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  updateProfile = (fieldValues) => {
    fieldValues.current_station_id = fieldValues.current_station.id;
    fieldValues.id = this.props.serviceMember.id;
    if (fieldValues.rank !== this.props.serviceMember.rank) {
      this.props.entitlementChanged();
    }
    fieldValues.move_id = this.props.move?.id;

    return patchServiceMember(fieldValues)
      .then((response) => {
        // Update Redux with new data
        this.props.updateServiceMember(response);

        this.props.editSuccessful();
        this.props.history.goBack();
        if (this.props.isPpm) {
          const moveId = this.props.move?.id;
          this.props.checkEntitlement(moveId);
        }
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');
        this.setState({
          errorMessage,
        });

        scrollToTop();
      });
  };

  componentDidMount() {
    this.props.editBegin();
    this.props.entitlementChangeBegin();
  }

  render() {
    const { schema, serviceMember, moveIsSubmitted, schemaAffiliation, schemaRank } = this.props;
    const { errorMessage } = this.state;

    return (
      <div className="usa-grid">
        {errorMessage && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {errorMessage}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditProfileForm
            initialValues={serviceMember}
            onSubmit={this.updateProfile}
            onCancel={this.returnToReview}
            schema={schema}
            moveIsSubmitted={moveIsSubmitted}
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
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    serviceMember,
    move: selectCurrentMove(state) || {},
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    moveIsSubmitted: selectMoveIsSubmitted(state),
    isPpm: selectHasCurrentPPM(state),
    schemaRank: get(state, 'swaggerInternal.spec.definitions.ServiceMemberRank', {}),
    schemaAffiliation: get(state, 'swaggerInternal.spec.definitions.Affiliation', {}),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      push,
      updateServiceMember: updateServiceMemberAction,
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
