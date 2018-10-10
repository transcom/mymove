import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, FormSection } from 'redux-form';

import { updateBackupInfo } from './ducks';

import { AddressElementDisplay, AddressElementEdit } from 'shared/Address';
import { validateRequiredFields } from 'shared/JsonSchemaForm';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelField, editablePanelify } from 'shared/EditablePanel';

const BackupInfoDisplay = props => {
  const backupAddress = props.backupMailingAddress;
  const backupContact = props.backupContact;

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <AddressElementDisplay
          address={backupAddress}
          title="Backup mailing address"
        />
      </div>
      <div className="editable-panel-column">
        <PanelField title="Backup contact">
          {backupContact.name}
          <br />
          {backupContact.telephone && (
            <span>
              {backupContact.telephone}
              <br />
            </span>
          )}
          {backupContact.email && (
            <span>
              {backupContact.email}
              <br />
            </span>
          )}
        </PanelField>
      </div>
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
          <SwaggerField fieldName="telephone" {...backupContactProps} />
          <SwaggerField fieldName="email" {...backupContactProps} required />

          <div className="panel-subhead">Authorization</div>
          <SwaggerField fieldName="permission" {...backupContactProps} />
        </FormSection>
      </div>

      <div className="editable-panel-column">
        <FormSection name="backupMailingAddress">
          <AddressElementEdit
            addressProps={backupMailingAddressProps}
            title="Backup mailing address"
          />
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_backup_info';

let BackupInfoPanel = editablePanelify(BackupInfoDisplay, BackupInfoEdit);
BackupInfoPanel = reduxForm({
  form: formName,
  validate: validateRequiredFields,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(BackupInfoPanel);

function mapStateToProps(state) {
  let serviceMember = get(state, 'office.officeServiceMember', {});
  let backupContact = get(state, 'office.officeBackupContacts.0', {}); // there can be only one

  return {
    // reduxForm
    initialValues: {
      backupContact: backupContact,
      backupMailingAddress: get(serviceMember, 'backup_mailing_address', {}),
    },

    addressSchema: get(state, 'swaggerInternal.spec.definitions.Address', {}),
    backupContactSchema: get(
      state,
      'swaggerInternal.spec.definitions.ServiceMemberBackupContactPayload',
      {},
    ),
    backupMailingAddress: get(serviceMember, 'backup_mailing_address', {}),
    backupContact: backupContact,

    // editablePanelify
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
