import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, FormSection } from 'redux-form';

import { AddressElementDisplay, AddressElementEdit } from 'shared/Address';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateRequiredFields } from 'shared/JsonSchemaForm';

import {
  updateServiceMember,
  loadServiceMemberLabel,
  updateServiceMemberLabel,
} from 'shared/Entities/modules/serviceMembers';
import { PanelSwaggerField, PanelField, SwaggerValue, editablePanelify } from 'shared/EditablePanel';
import { stringifyName } from 'shared/utils/serviceMember';
import { getRequestStatus } from 'shared/Swagger/selectors';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';

const CustomerInfoDisplay = props => {
  const fieldProps = {
    schema: props.serviceMemberSchema,
    values: props.serviceMember,
  };
  const values = props.serviceMember;
  const name = stringifyName(values);
  const address = get(values, 'residential_address', {});

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelField title="Name" value={name} />
        <PanelSwaggerField title="DoD ID" fieldName="edipi" {...fieldProps} />
        <PanelField title="Branch & rank">
          <SwaggerValue fieldName="affiliation" {...fieldProps} /> - <SwaggerValue fieldName="rank" {...fieldProps} />
        </PanelField>
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField title="Phone" fieldName="telephone" {...fieldProps} />
        <PanelSwaggerField title="Alt. Phone" fieldName="secondary_telephone" {...fieldProps} />
        <PanelSwaggerField title="Email" fieldName="personal_email" {...fieldProps} />
        <PanelField title="Pref. contact" className="contact-prefs">
          {values.phone_is_preferred && (
            <span>
              <FontAwesomeIcon icon={faPhone} flip="horizontal" />
              phone
            </span>
          )}
          {values.text_message_is_preferred && (
            <span>
              <FontAwesomeIcon icon={faComments} />
              text
            </span>
          )}
          {values.email_is_preferred && (
            <span>
              <FontAwesomeIcon icon={faEmail} />
              email
            </span>
          )}
        </PanelField>
        <AddressElementDisplay address={address} title="Current Address" />
      </div>
    </React.Fragment>
  );
};

const CustomerInfoEdit = props => {
  const schema = props.serviceMemberSchema;
  let addressProps = {
    swagger: props.addressSchema,
    values: props.initialValues.address,
  };

  return (
    <React.Fragment>
      <div>
        <FormSection name="serviceMember">
          <div className="editable-panel-column">
            <SwaggerField fieldName="title" swagger={schema} />
            <SwaggerField fieldName="first_name" swagger={schema} required />
            <SwaggerField fieldName="middle_name" swagger={schema} />
            <SwaggerField fieldName="last_name" swagger={schema} required />
            <SwaggerField fieldName="suffix" swagger={schema} />
          </div>
          <div className="editable-panel-column">
            <SwaggerField fieldName="edipi" swagger={schema} required />
            <SwaggerField fieldName="affiliation" swagger={schema} required />
            <SwaggerField fieldName="rank" swagger={schema} required />
          </div>
        </FormSection>
      </div>

      <div>
        <div className="editable-panel-column">
          <FormSection name="serviceMember">
            <div className="panel-subhead">Contact</div>
            <SwaggerField fieldName="telephone" swagger={schema} required />
            <SwaggerField fieldName="secondary_telephone" swagger={schema} />
            <SwaggerField fieldName="personal_email" swagger={schema} required />

            <fieldset key="contact_preferences">
              <legend htmlFor="contact_preferences">
                <p>Preferred contact method</p>
              </legend>
              <SwaggerField fieldName="phone_is_preferred" swagger={schema} />
              <SwaggerField fieldName="text_message_is_preferred" swagger={schema} />
              <SwaggerField fieldName="email_is_preferred" swagger={schema} />
            </fieldset>
          </FormSection>
        </div>

        <div className="editable-panel-column">
          <FormSection name="address">
            <AddressElementEdit addressProps={addressProps} title="Current Residence Address" />
          </FormSection>
        </div>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_customer_info';

let CustomerInfoPanel = editablePanelify(CustomerInfoDisplay, CustomerInfoEdit);
CustomerInfoPanel = reduxForm({
  form: formName,
  validate: validateRequiredFields,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(CustomerInfoPanel);

function mapStateToProps(state, ownProps) {
  let customerInfo = ownProps.serviceMember;
  const loadServiceMemberStatus = getRequestStatus(state, loadServiceMemberLabel);
  const updateServiceMemberStatus = getRequestStatus(state, updateServiceMemberLabel);
  let hasError = false;
  let errorMessage = '';

  if (loadServiceMemberStatus.isSuccess === false) {
    hasError = true;
    errorMessage = get(loadServiceMemberStatus, 'error.response.message', '');
  }

  if (updateServiceMemberStatus.isSuccess === false) {
    hasError = true;
    errorMessage = get(updateServiceMemberStatus, 'error.response.message', '');
  }

  return {
    // reduxForm
    initialValues: {
      serviceMember: customerInfo,
      address: customerInfo.residential_address,
    },

    addressSchema: get(state, 'swaggerInternal.spec.definitions.Address', {}),

    // CustomerInfoEdit
    serviceMemberSchema: get(state, 'swaggerInternal.spec.definitions.ServiceMemberPayload'),
    serviceMember: ownProps.serviceMember,

    hasError,
    errorMessage,
    isUpdating: false,

    // editablePanelify
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      let serviceMember = values.serviceMember;
      serviceMember.residential_address = values.address;
      return [ownProps.serviceMember.id, serviceMember];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateServiceMember,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(CustomerInfoPanel);
