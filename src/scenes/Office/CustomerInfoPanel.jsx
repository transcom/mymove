// import { get, pick } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import editablePanel from './editablePanel';

import { no_op_action } from 'shared/utils';

// import { updateCustomerInfo, loadCustomerInfo } from './ducks';
// import { PanelField } from 'shared/EditablePanel';
// import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const CustomerInfoDisplay = props => {
  // const fieldProps = pick(props, ['schema', 'values']);
  return (
    <React.Fragment>
      <div className="editable-panel-column" />
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
    schema: {},
    hasError: false,
    errorMessage: state.office.error,
    displayValues: {},
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
