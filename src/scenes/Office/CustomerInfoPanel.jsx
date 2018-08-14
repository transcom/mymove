import { get, compact } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from './editablePanel';
import { addressElementDisplay, addressElementEdit } from './AddressElement';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateRequiredFields } from 'shared/JsonSchemaForm';

import { updateServiceMember } from './ducks';
import {
  PanelSwaggerField,
  PanelField,
  SwaggerValue,
} from 'shared/EditablePanel';

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
  const name = compact([values.last_name, values.first_name]).join(', ');
  const address = get(values, 'residential_address', {});

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelField title="Name" value={name} />
        <PanelSwaggerField title="DoD ID" fieldName="edipi" {...fieldProps} />
        <PanelField title="Branch & rank">
          <SwaggerValue fieldName="affiliation" {...fieldProps} /> -{' '}
          <SwaggerValue fieldName="rank" {...fieldProps} />
        </PanelField>
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField
          title="Phone"
          fieldName="telephone"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Alt. Phone"
          fieldName="secondary_telephone"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Email"
          fieldName="personal_email"
          {...fieldProps}
        />
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
        {addressElementDisplay(address, 'Current Address')}
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
            <SwaggerField
              fieldName="personal_email"
              swagger={schema}
              required
            />

            <fieldset key="contact_preferences">
              <legend htmlFor="contact_preferences">
                <p>Preferred contact method</p>
              </legend>
              <SwaggerField fieldName="phone_is_preferred" swagger={schema} />
              <SwaggerField
                fieldName="text_message_is_preferred"
                swagger={schema}
              />
              <SwaggerField fieldName="email_is_preferred" swagger={schema} />
            </fieldset>
          </FormSection>
        </div>

        <div className="editable-panel-column">
          <FormSection name="address">
            {addressElementEdit(addressProps, 'Current Residence Address')}
          </FormSection>
        </div>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_customer_info';

let CustomerInfoPanel = editablePanel(CustomerInfoDisplay, CustomerInfoEdit);
CustomerInfoPanel = reduxForm({
  form: formName,
  validate: validateRequiredFields,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(CustomerInfoPanel);

function mapStateToProps(state) {
  let customerInfo = get(state, 'office.officeServiceMember', {});
  return {
    // reduxForm
    initialValues: {
      serviceMember: customerInfo,
      address: customerInfo.residential_address,
    },

    addressSchema: get(state, 'swagger.spec.definitions.Address', {}),

    // CustomerInfoEdit
    serviceMemberSchema: get(
      state,
      'swagger.spec.definitions.ServiceMemberPayload',
    ),
    serviceMember: state.office.officeServiceMember,

    hasError:
      state.office.serviceMemberHasLoadError ||
      state.office.serviceMemberHasUpdateError,
    errorMessage: state.office.error,
    isUpdating: false,

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      let serviceMember = values.serviceMember;
      serviceMember.residential_address = values.address;
      return [state.office.officeServiceMember.id, serviceMember];
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
