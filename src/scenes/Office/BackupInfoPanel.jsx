import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, FormSection } from 'redux-form';
import editablePanel from './editablePanel';

import { updateBackupInfo } from './ducks';

// import { PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const BackupInfoDisplay = props => {
  //const fieldProps = pick(props, ['schema', 'values']);
  return (
    <React.Fragment>
      <div className="editable-panel-column" />
    </React.Fragment>
  );
};

const BackupInfoEdit = props => {
  let backupContactProps = {
    swagger: props.backupContactSchema,
    values: props.backupContact,
  };
  let backupMailingAddressProps = {
    swagger: props.addressSchema,
    values: props.backupMailingAddress,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <div className="panel-subhead">Backup Contact 1</div>

        <FormSection name="backupContact">
          <SwaggerField fieldName="name" {...backupContactProps} required />
          <SwaggerField
            fieldName="telephone"
            {...backupContactProps}
            required
          />
          <SwaggerField fieldName="email" {...backupContactProps} />

          <div className="panel-subhead">Authorization</div>
          <SwaggerField fieldName="permission" {...backupContactProps} />
        </FormSection>
      </div>

      <div className="editable-panel-column">
        <div className="panel-subhead">Backup Mailing Address</div>

        <FormSection name="backupMailingAddress">
          <SwaggerField
            fieldName="street_address_1"
            {...backupMailingAddressProps}
          />
          <SwaggerField
            fieldName="street_address_2"
            {...backupMailingAddressProps}
          />
          <SwaggerField fieldName="city" {...backupMailingAddressProps} />
          <SwaggerField fieldName="state" {...backupMailingAddressProps} />
          <SwaggerField
            fieldName="postal_code"
            {...backupMailingAddressProps}
          />
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_backup_info';

let BackupInfoPanel = editablePanel(BackupInfoDisplay, BackupInfoEdit);
BackupInfoPanel = reduxForm({ form: formName })(BackupInfoPanel);

function mapStateToProps(state) {
  let serviceMember = get(state, 'office.officeServiceMember', {});
  let backupContact = get(state, 'office.officeBackupContacts.0', {}); // there can be only one

  return {
    // reduxForm
    formData: state.form[formName],
    initialValues: {
      backupContact: backupContact,
      backupMailingAddress: serviceMember.backup_mailing_address,
    },

    addressSchema: get(state, 'swagger.spec.definitions.Address', {}),
    backupContactSchema: get(
      state,
      'swagger.spec.definitions.ServiceMemberBackupContactPayload',
      {},
    ),
    backupMailingAddress: serviceMember.backup_mailing_address,
    backupContact: backupContact,

    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      return [
        serviceMember.id,
        { backup_mailing_address: values.backupMailingAddress },
        backupContact.id,
        values.backupContact,
      ];
    },

    hasError: false,
    errorMessage: state.office.error,
    isUpdating: false,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateBackupInfo,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(BackupInfoPanel);
