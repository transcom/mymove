import { get, compact } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import editablePanel from './editablePanel';

import { no_op_action } from 'shared/utils';

// import { updateCustomerInfo, loadCustomerInfo } from './ducks';
import {
  PanelSwaggerField,
  PanelField,
  SwaggerValue,
} from 'shared/EditablePanel';
// import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';

const CustomerInfoDisplay = props => {
  const serviceMemberFieldProps = {
    schema: props.serviceMemberSchema,
    values: props.displayServiceMemberValues,
  };
  const values = props.displayServiceMemberValues;
  const name = compact([values.last_name, values.first_name]).join(', ');
  const address = values.residential_address || {};

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelField title="Name" value={name} />
        <PanelSwaggerField
          title="DoD ID"
          fieldName="edipi"
          {...serviceMemberFieldProps}
        />
        <PanelField title="Branch & rank">
          <SwaggerValue fieldName="affiliation" {...serviceMemberFieldProps} />{' '}
          - <SwaggerValue fieldName="rank" {...serviceMemberFieldProps} />
        </PanelField>
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField
          title="Phone"
          fieldName="telephone"
          {...serviceMemberFieldProps}
        />
        <PanelSwaggerField
          title="Alt. Phone"
          fieldName="secondary_telephone"
          {...serviceMemberFieldProps}
        />
        <PanelSwaggerField
          title="Email"
          fieldName="personal_email"
          {...serviceMemberFieldProps}
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
  // const { schema } = props;
  return (
    <React.Fragment>
      <div className="form-column">
        <label>Title (optional)</label>
        <input type="text" name="title" />
      </div>
      <div className="form-column">
        <label>First name</label>
        <input type="text" name="first-name" />
      </div>
      <div className="form-column">
        <label>Middle name (optional)</label>
        <input type="text" name="middle-name" />
      </div>
      <div className="form-column">
        <label>Last name</label>
        <input type="text" name="last-name" />
      </div>
      <div className="form-column">
        <label>Suffix (optional)</label>
        <input type="text" name="name-suffix" />
      </div>
      <div className="form-column">
        <label>DoD ID</label>
        <input type="number" name="dod-id" />
      </div>
      <div className="form-column">
        <label>Branch</label>
        <select name="branch">
          <option value="army">Army</option>
          <option value="navy">Navy</option>
          <option value="air-force">Air Force</option>
          <option value="marines">Marines</option>
          <option value="coast-guard">Coast Guard</option>
        </select>
      </div>
      <div className="form-column">
        <label>Rank</label>
        <select name="rank">
          <option value="E-7">E-7</option>
          <option value="another-rank">Another rank</option>
          <option value="and-another-rank">And another rank</option>
        </select>
      </div>
      <div className="form-column">
        <b>Contact</b>
        <label>Phone</label>
        <input type="tel" name="contact-phone-number" />
      </div>
      <div className="form-column">
        <label>Alternate phone</label>
        <input type="tel" name="alternate-contact-phone-number" />
      </div>
      <div className="form-column">
        <label>Email</label>
        <input type="text" name="contact-email" />
      </div>
      <div className="form-column">
        <label>Preferred contact methods</label>
        <div>
          <input
            type="checkbox"
            id="phone-preference"
            name="preferred-contact-phone"
          />
          <label htmlFor="phone-preference">Phone</label>
        </div>
        <div>
          <input
            type="checkbox"
            id="text-preference"
            name="preferred-contact-text-message"
          />
          <label htmlFor="text-preference">Text message</label>
        </div>
        <div>
          <input
            type="checkbox"
            id="email-preference"
            name="preferred-contact-email"
          />
          <label htmlFor="email-preference">Email</label>
        </div>
      </div>
      <div className="form-column">
        <b>Current Residence Address</b>
        <label>Address 1</label>
        <input type="text" name="contact-address-1" />
      </div>
      <div className="form-column">
        <label>Address 2</label>
        <input type="text" name="contact-address-2" />
      </div>
      <div className="form-column">
        <label>City</label>
        <input type="text" name="contact-city" />
      </div>
      <div className="form-column">
        <label>State</label>
        <input type="text" name="contact-state" />
      </div>
      <div className="form-column">
        <label>Zip</label>
        <input type="number" name="contact-zip" />
      </div>
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
    serviceMemberSchema: get(
      state,
      'swagger.spec.definitions.ServiceMemberPayload',
    ),
    hasError: false,
    errorMessage: state.office.error,
    displayServiceMemberValues: state.office.officeServiceMember,
    isUpdating: false,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op_action,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(CustomerInfoPanel);
