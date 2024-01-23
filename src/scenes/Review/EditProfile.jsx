import React, { Component } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { Field, reduxForm } from 'redux-form';

import SaveCancelButtons from './SaveCancelButtons';
import profileImage from './images/profile.png';

import { getResponseError, patchOrders, patchServiceMember } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { updateOrders as updateOrderAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';
import LocationSearchBox from 'components/LocationSearchBox/LocationSearchBox';
import scrollToTop from 'shared/scrollToTop';
import {
  selectCurrentMove,
  selectCurrentOrders,
  selectWeightAllotmentsForLoggedInUser,
  selectHasCurrentPPM,
  selectMoveIsInDraft,
  selectServiceMemberFromLoggedInUser,
} from 'store/entities/selectors';

import './Review.css';

import SectionWrapper from 'components/Customer/SectionWrapper';
import ServiceInfoDisplay from 'components/Customer/Review/ServiceInfoDisplay/ServiceInfoDisplay';

const editProfileFormName = 'edit_profile';

let EditProfileForm = (props) => {
  const { schema, handleSubmit, submitting, valid, moveIsInDraft, initialValues, serviceMember } = props;
  const currentDutyLocation = get(serviceMember, 'current_location');
  const transportationOfficeName = get(currentDutyLocation, 'transportation_office.name');
  const transportationOfficePhone = get(currentDutyLocation, 'transportation_office.phone_lines.0');
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
              {moveIsInDraft && (
                <>
                  <SwaggerField fieldName="affiliation" swagger={schema} required />
                  <SwaggerField fieldName="rank" swagger={schema} required />
                  <SwaggerField fieldName="edipi" swagger={schema} required />
                  <Field name="current_location" title="Current duty location" component={LocationSearchBox} />
                </>
              )}
              {!moveIsInDraft && (
                <ServiceInfoDisplay
                  firstName={initialValues.first_name}
                  lastName={initialValues.last_name}
                  originDutyLocationName={currentDutyLocation.name}
                  originTransportationOfficeName={transportationOfficeName}
                  originTransportationOfficePhone={transportationOfficePhone}
                  affiliation={initialValues.affiliation}
                  payGrade={initialValues.rank}
                  edipi={initialValues.edipi}
                  isEditable={false}
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
const validateProfileForm = validateAdditionalFields(['current_location']);
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
    const { setFlashMessage, entitlement } = this.props;

    let entitlementCouldChange = false;

    fieldValues.current_location_id = fieldValues.current_location.id;
    fieldValues.id = this.props.serviceMember.id;

    if (fieldValues.rank !== this.props.currentOrders.grade) {
      entitlementCouldChange = true;
    }

    patchServiceMember(fieldValues)
      .then((response) => {
        // Update Redux with new data
        this.props.updateServiceMember(response);

        if (entitlementCouldChange) {
          setFlashMessage(
            'EDIT_PROFILE_SUCCESS',
            'info',
            `Your weight entitlement is now ${entitlement.sum.toLocaleString()} lbs.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_PROFILE_SUCCESS', 'success', '', 'Your changes have been saved.');
        }

        const { router: navigate } = this.props;
        navigate(-1);
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

    patchOrders(fieldValues)
      .then((response) => {
        // Update Redux with new data
        this.props.updateOrders(response);

        if (entitlementCouldChange) {
          setFlashMessage(
            'EDIT_PROFILE_SUCCESS',
            'info',
            `Your weight entitlement is now ${entitlement.sum.toLocaleString()} lbs.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_PROFILE_SUCCESS', 'success', '', 'Your changes have been saved.');
        }

        const { router: navigate } = this.props;
        navigate(-1);
      })
      .catch((e) => {
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update orders due to server error');
        this.setState({
          errorMessage,
        });

        scrollToTop();
      });
  };

  render() {
    const { schema, serviceMember, moveIsInDraft, schemaAffiliation, schemaRank, currentOrders } = this.props;
    const { errorMessage } = this.state;
    const initialValues = {
      ...serviceMember,
      payGrade: currentOrders.grade,
      current_location: currentOrders.origin_duty_location,
    };
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
            initialValues={initialValues}
            onSubmit={this.updateProfile}
            onCancel={this.returnToReview}
            schema={schema}
            moveIsInDraft={moveIsInDraft}
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
    currentOrders: selectCurrentOrders(state),
    // The move still counts as in draft if there are no orders.
    moveIsInDraft: selectMoveIsInDraft(state) || !selectCurrentOrders(state),
    isPpm: selectHasCurrentPPM(state),
    schemaGrade: get(state, 'swaggerInternal.spec.definitions.OrderPayGrade', {}),
    schemaAffiliation: get(state, 'swaggerInternal.spec.definitions.Affiliation', {}),
    entitlement: selectWeightAllotmentsForLoggedInUser(state),
  };
}

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
  updateOrders: updateOrderAction,
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(EditProfile);
