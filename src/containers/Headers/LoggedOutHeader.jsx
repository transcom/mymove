import React, { useState } from 'react';

import ConnectedEulaModal from 'components/EulaModal';
import MilMoveHeader from 'components/MilMoveHeader/index';
import LoggedOutUserInfo from 'components/MilMoveHeader/LoggedOutUserInfo';

const LoggedOutHeader = () => {
  const [showEula, setShowEula] = useState(false);

  return (
    <>
      <ConnectedEulaModal
        isOpen={showEula}
        acceptTerms={() => {
          window.location.href = '/auth/okta';
        }}
        closeModal={() => setShowEula(false)}
      />

      <MilMoveHeader>
        <LoggedOutUserInfo handleLogin={() => setShowEula(true)} />
      </MilMoveHeader>
    </>
  );
};

export default LoggedOutHeader;
