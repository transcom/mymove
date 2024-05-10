import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { generalRoutes } from 'constants/routes';
import ConnectedEulaModal from 'components/EulaModal';
import MilMoveHeader from 'components/MilMoveHeader/index';
import LoggedOutUserInfo from 'components/MilMoveHeader/LoggedOutUserInfo';

const LoggedOutHeader = ({ app }) => {
  const [showEula, setShowEula] = useState(false);

  const navigate = useNavigate();

  const handleRequestAccount = () => {
    navigate(generalRoutes.REQUEST_ACCOUNT);
  };

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
        <LoggedOutUserInfo
          handleLogin={() => setShowEula(true)}
          handleRequestAccount={() => handleRequestAccount()}
          app={app}
        />
      </MilMoveHeader>
    </>
  );
};

export default LoggedOutHeader;
