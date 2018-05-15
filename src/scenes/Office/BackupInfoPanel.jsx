//import { pick } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import editablePanel from './editablePanel';

import { no_op_action } from 'shared/utils';

// import { updateBackupInfo, loadBackupInfo } from './ducks';
// import { PanelField } from 'shared/EditablePanel';
// import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const BackupInfoDisplay = props => {
  //const fieldProps = pick(props, ['schema', 'values']);
  return (
    <React.Fragment>
      <div className="editable-panel-column" />
    </React.Fragment>
  );
};

const BackupInfoEdit = props => {
  // const { schema } = props;
  return (
    <React.Fragment>
      <div className="form-column">
        <b>Backup Contact 1</b>
        <label>Name</label>
        <input type="text" name="backup-contact-1-name" />
      </div>
      <div className="form-column">
        <label>Phone</label>
        <input type="tel" name="backup-contact-1-phone" />
      </div>
      <div className="form-column">
        <label>Email (optional)</label>
        <input type="text" name="backup-contact-1-email" />
      </div>
      <div className="form-column">
        <b>Authorization</b>
        <input type="radio" name="authorization" value="none" />
        <label htmlFor="none">None</label>
        <input
          type="radio"
          name="authorization"
          value="letter-of-authorization"
        />
        <label htmlFor="letter-of-authorization">Letter of Authorization</label>
        <input
          type="radio"
          name="authorization"
          value="sign-for-pickup-delivery-only"
        />
        <label htmlFor="sign-for-pickup-delivery-only">
          Sign for pickup/delivery only
        </label>
      </div>
      <div className="form-column">
        <b>Backup Mailing Address</b>
        <label>Mailing address 1</label>
        <input type="text" name="backup-contact-1-mailing-address-1" />
      </div>
      <div className="form-column">
        <label>Mailing address 2</label>
        <input type="text" name="backup-contact-1-mailing-address-2" />
      </div>
      <div className="form-column">
        <label>City</label>
        <input type="text" name="backup-contact-1-city" />
      </div>
      <div className="form-column">
        <label>State</label>
        <input type="text" name="backup-contact-1-state" />
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_backup_info';

let BackupInfoPanel = editablePanel(BackupInfoDisplay, BackupInfoEdit);
BackupInfoPanel = reduxForm({ form: formName })(BackupInfoPanel);

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

export default connect(mapStateToProps, mapDispatchToProps)(BackupInfoPanel);
