import React from 'react';
import { Confirm } from 'react-admin';

import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

const RequestedOfficeUserPrivilegeConfirm = ({
  dialogId,
  isOpen,
  title = 'Attention: The user has requested the selected privilege(s)',
  privileges = [],
  checkedPrivileges = [],
  setCheckedPrivileges,
  onConfirm,
  onClose,
}) => {
  const filteredPrivileges = (privileges || []).filter(
    (priv) => priv.privilegeType === elevatedPrivilegeTypes.SUPERVISOR,
  );

  return (
    <Confirm
      isOpen={isOpen}
      title={title}
      content={
        <div id={dialogId} data-testid="RequestedOfficeUserPrivilegeConfirm">
          <p id="privilege-dialog-desc" aria-labelledby="privilege-dialog-legend">
            If the user is not qualified for a selected privilege, please deselect it before approval.
            <br />
            If you want to halt the approval process, select Cancel.
          </p>
          <fieldset
            aria-labelledby="privilege-dialog-legend privilege-dialog-desc"
            aria-describedby="privilege-dialog-desc"
            style={{ margin: '1rem 0', border: 0, padding: 0 }}
          >
            <legend id="privilege-dialog-legend" className="usa-sr-only">
              Requested privileges
            </legend>
            {filteredPrivileges.length > 0 && (
              <>
                {filteredPrivileges.map((priv) => (
                  <div key={priv.id} style={{ display: 'flex', alignItems: 'center', marginBottom: 4 }}>
                    <input
                      type="checkbox"
                      id={`privilege-${priv.id}`}
                      name="privileges"
                      value={priv.id}
                      checked={checkedPrivileges.includes(priv.id)}
                      aria-labelledby={`privilege-label-${priv.id}`}
                      aria-describedby="privilege-dialog-desc"
                      tabIndex={0}
                      onKeyDown={(e) => {
                        if (e.key === ' ' || e.key === 'Enter') {
                          e.preventDefault();
                          setCheckedPrivileges((prev) =>
                            prev.includes(priv.id) ? prev.filter((id) => id !== priv.id) : [...prev, priv.id],
                          );
                        }
                      }}
                      onChange={(e) => {
                        setCheckedPrivileges((prev) =>
                          e.target.checked ? [...prev, priv.id] : prev.filter((id) => id !== priv.id),
                        );
                      }}
                    />
                    <label id={`privilege-label-${priv.id}`} htmlFor={`privilege-${priv.id}`} style={{ marginLeft: 8 }}>
                      {priv.privilegeName && String(priv.privilegeName).trim().length > 0
                        ? priv.privilegeName
                        : 'Supervisor'}
                    </label>
                  </div>
                ))}
              </>
            )}
          </fieldset>
        </div>
      }
      onConfirm={onConfirm}
      onClose={onClose}
    />
  );
};

export default RequestedOfficeUserPrivilegeConfirm;
