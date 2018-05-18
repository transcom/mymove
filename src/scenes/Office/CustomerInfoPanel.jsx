import { get, pick, compact } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';
import editablePanel from './editablePanel';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { no_op_action } from 'shared/utils';

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
  const fieldProps = pick(props, ['schema', 'values']);
  const values = props.values;
  const name = compact([values.last_name, values.first_name]).join(', ');
  const address = values.residential_address || {};
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
        <PanelField title="Current Address">
          {address.street_address_1}
          <br />
          {address.street_address_2 && (
            <span>
              {address.street_address_2}
              <br />
            </span>
          )}
          {address.street_address_3 && (
            <span>
              {address.street_address_3}
              <br />
            </span>
          )}
          {address.city}, {address.state} {address.postal_code}
        </PanelField>
      </div>
    </React.Fragment>
  );
};

const CustomerInfoEdit = props => {
  const schema = props.schema;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <SwaggerField fieldName="title" swagger={schema} />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="first_name" swagger={schema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="middle_name" swagger={schema} />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="last_name" swagger={schema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="suffix" swagger={schema} />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="edipi" swagger={schema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="affiliation" swagger={schema} />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="rank" swagger={schema} />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="telephone" swagger={schema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="secondary_telephone" swagger={schema} />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="personal_email" swagger={schema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField fieldName="residential_address" swagger={schema} />
      </div>
      <fieldset key="contact_preferences">
        <legend htmlFor="contact_preferences">
          Preferred contact method during your move:
        </legend>
        <SwaggerField fieldName="phone_is_preferred" swagger={schema} />
        <SwaggerField fieldName="text_message_is_preferred" swagger={schema} />
        <SwaggerField fieldName="email_is_preferred" swagger={schema} />
      </fieldset>
    </React.Fragment>
  );
};

const formName = 'office_move_info_customer_info';

let CustomerInfoPanel = editablePanel(CustomerInfoDisplay, CustomerInfoEdit);
CustomerInfoPanel = reduxForm({ form: formName })(CustomerInfoPanel);

function mapStateToProps(state) {
  return {
    // reduxForm
    formData: state.form[formName],
    initialValues: {},

    // Wrapper
    schema: get(state, 'swagger.spec.definitions.ServiceMemberPayload'),
    hasError: false,
    errorMessage: state.office.error,
    displayValues: state.office.officeServiceMember || {},
    isUpdating: false,
    // formValues: getFormValues(formName)(state),
    getUpdateArgs: function() {
      return [
        state.office.officeServiceMember.id,
        getFormValues(formName)(state),
      ];
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
