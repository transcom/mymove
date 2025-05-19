// import React, { useEffect, useState } from 'react';
// import { generatePath, useNavigate, useParams } from 'react-router-dom';
// import { useDispatch, useSelector } from 'react-redux';
// import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

// import NotificationScrollToTop from 'components/NotificationScrollToTop';
// import {
//   selectMTOShipmentById,
//   selectGunsafeWeightTicketAndIndexById,
//   selectServiceMemberFromLoggedInUser,
// } from 'store/entities/selectors';
// import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
// import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
// import { shipmentTypes } from 'constants/shipments';
// import { customerRoutes } from 'constants/routes';
// import {
//   createUploadForPPMDocument,
//   createGunsafeWeightTicket,
//   deleteUpload,
//   patchGunsafeWeightTicket,
//   patchMTOShipment,
//   getMTOShipmentsForMove,
//   getAllMoves,
// } from 'services/internalApi';
// import LoadingPlaceholder from 'shared/LoadingPlaceholder';
// import GunsafeForm from 'components/Shared/PPM/Closeout/GunsafeForm/GunsafeForm';
// import { updateAllMoves, updateMTOShipment } from 'store/entities/actions';
// import { CUSTOMER_ERROR_MESSAGES } from 'constants/errorMessages';
// import { APP_NAME } from 'constants/apps';

// const Gunsafe = () => {
//   const dispatch = useDispatch();
//   const navigate = useNavigate();

//   const serviceMember = useSelector((state) => selectServiceMemberFromLoggedInUser(state));
//   const serviceMemberId = serviceMember.id;

//   const gunsafeEntitlements = useSelector((state) => selectGunsafeEntitlements(state));

//   const appName = APP_NAME.MYMOVE;
//   const { moveId, mtoShipmentId, gunsafeId } = useParams();

//   const handleBack = () => {
//     const path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
//       moveId,
//       mtoShipmentId,
//     });
//     navigate(path);
//   };
//   const [errorMessage, setErrorMessage] = useState(null);

//   const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
//   const { gunsafeWeightTicket: currentGunsafeWeightTicket, index: currentIndex } = useSelector((state) =>
//     selectGunsafeWeightTicketAndIndexById(state, mtoShipmentId, gunsafeId),
//   );

//   useEffect(() => {
//     if (!gunsafeId) {
//       createGunsafeWeightTicket(mtoShipment?.ppmShipment?.id)
//         .then((resp) => {
//           if (mtoShipment?.ppmShipment?.gunsafeWeightTickets) {
//             mtoShipment.ppmShipment.gunsafeWeightTickets.push(resp);
//           } else {
//             mtoShipment.ppmShipment.gunsafeWeightTickets = [resp];
//           }
//           // Update the URL so the back button would work and not create a new weight ticket or on
//           // refresh either.
//           navigate(
//             generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
//               moveId,
//               mtoShipmentId,
//               gunsafeId: resp.id,
//             }),
//             { replace: true },
//           );
//           dispatch(updateMTOShipment(mtoShipment));
//         })
//         .catch(() => {
//           setErrorMessage('Failed to create trip record');
//         });
//     }
//   }, [gunsafeId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment]);

//   useEffect(() => {
//     const moves = getAllMoves(serviceMemberId);
//     dispatch(updateAllMoves(moves));
//   }, [gunsafeId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment, serviceMemberId]);

//   const handleErrorMessage = (error) => {
//     if (error?.response?.status === 412) {
//       setErrorMessage(CUSTOMER_ERROR_MESSAGES.PRECONDITION_FAILED);
//     } else {
//       setErrorMessage('Failed to fetch shipment information');
//     }
//   };

//   const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
//     const documentId = currentGunsafeWeightTicket[`${fieldName}Id`];

//     // Create a date-time stamp in the format "yyyymmddhh24miss"
//     const now = new Date();
//     const timestamp =
//       now.getFullYear().toString() +
//       (now.getMonth() + 1).toString().padStart(2, '0') +
//       now.getDate().toString().padStart(2, '0') +
//       now.getHours().toString().padStart(2, '0') +
//       now.getMinutes().toString().padStart(2, '0') +
//       now.getSeconds().toString().padStart(2, '0');

//     // Create a new filename with the timestamp prepended
//     const newFileName = `${file.name}-${timestamp}`;

//     // Create and return a new File object with the new filename
//     const newFile = new File([file], newFileName, { type: file.type });

//     createUploadForPPMDocument(mtoShipment.ppmShipment.id, documentId, newFile, false)
//       .then((upload) => {
//         mtoShipment.ppmShipment.gunsafeWeightTickets[currentIndex][fieldName].uploads.push(upload);
//         dispatch(updateMTOShipment(mtoShipment));
//         setFieldTouched(fieldName, true);
//         return upload;
//       })
//       .catch(() => {
//         setErrorMessage('Failed to save the file upload');
//       });
//   };

//   const handleUploadComplete = (err) => {
//     if (err) {
//       setErrorMessage('Encountered error when completing file upload');
//     }
//   };

//   const handleUploadDelete = (uploadId, fieldName, setFieldTouched, setFieldValue) => {
//     deleteUpload(uploadId, null, mtoShipment?.ppmShipment?.id)
//       .then(() => {
//         const filteredUploads = mtoShipment.ppmShipment.gunsafeWeightTickets[currentIndex][fieldName].uploads.filter(
//           (upload) => upload.id !== uploadId,
//         );
//         mtoShipment.ppmShipment.gunsafeWeightTickets[currentIndex][fieldName].uploads = filteredUploads;
//         setFieldValue(fieldName, filteredUploads, true);
//         setFieldTouched(fieldName, true, true);
//         dispatch(updateMTOShipment(mtoShipment));
//       })
//       .catch(() => {
//         setErrorMessage('Failed to delete the file upload');
//       });
//   };

//   const updateMtoShipment = (values) => {
//     const belongsToSelf = values.belongsToSelf === 'true';
//     let gunsafe;
//     let spouseGunsafe;
//     if (belongsToSelf) {
//       gunsafe = values.weight;
//     }
//     if (!belongsToSelf) {
//       spouseGunsafe = values.weight;
//     }
//     const payload = {
//       belongsToSelf,
//       ppmShipment: {
//         id: mtoShipment.ppmShipment.id,
//       },
//       shipmentType: mtoShipment.shipmentType,
//       actualSpouseGunsafeWeight: parseInt(spouseGunsafe, 10),
//       actualGunsafeWeight: parseInt(gunsafe, 10),
//       shipmentLocator: values.shipmentLocator,
//       eTag: mtoShipment.eTag,
//     };

//     patchMTOShipment(mtoShipment.id, payload, payload.eTag)
//       .then((response) => {
//         navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
//         dispatch(updateMTOShipment(response));
//       })
//       .catch(() => {
//         setErrorMessage('Failed to update MTO shipment due to server error.');
//       });
//   };

//   const updateGunsafeWeightTicket = (values) => {
//     const hasWeightTickets = !values.missingWeightTicket;
//     const belongsToSelf = values.belongsToSelf === 'true';
//     const payload = {
//       ppmShipmentId: mtoShipment.ppmShipment.id,
//       gunsafeWeightTicketId: currentGunsafeWeightTicket.id,
//       description: values.description,
//       weight: parseInt(values.weight, 10),
//       belongsToSelf,
//       hasWeightTickets,
//     };

//     patchGunsafeWeightTicket(
//       mtoShipment?.ppmShipment?.id,
//       currentGunsafeWeightTicket.id,
//       payload,
//       currentGunsafeWeightTicket.eTag,
//     )
//       .then((resp) => {
//         mtoShipment.ppmShipment.gunsafeWeightTickets[currentIndex] = resp;
//         getMTOShipmentsForMove(moveId)
//           .then((response) => {
//             dispatch(updateMTOShipment(response.mtoShipments[mtoShipmentId]));
//             mtoShipment.eTag = response.mtoShipments[mtoShipmentId].eTag;
//             updateMtoShipment(values);
//             navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
//           })
//           .catch(() => {
//             setErrorMessage('Failed to fetch shipment information');
//           });
//       })
//       .catch((error) => {
//         handleErrorMessage(error);
//       });
//   };

//   const handleSubmit = async (values, { setSubmitting, setErrors }) => {
//     setErrorMessage(null);
//     setErrors({});
//     setSubmitting(false);
//     updateGunsafeWeightTicket(values);
//   };

//   const renderError = () => {
//     if (!errorMessage) {
//       return null;
//     }

//     return (
//       <Alert slim type="error">
//         {errorMessage}
//       </Alert>
//     );
//   };

//   if (!mtoShipment || !currentGunsafeWeightTicket) {
//     return renderError() || <LoadingPlaceholder />;
//   }
//   return (
//     <div className={ppmPageStyles.ppmPageStyle}>
//       <NotificationScrollToTop dependency={errorMessage} />
//       <GridContainer>
//         <Grid row>
//           <Grid col desktop={{ col: 8, offset: 2 }}>
//             <ShipmentTag shipmentType={shipmentTypes.PPM} />
//             <h1>Pro-gear</h1>
//             {renderError()}
//             <GunsafeForm
//               entitlements={gunsafeEntitlements}
//               gunsafe={currentGunsafeWeightTicket}
//               setNumber={currentIndex + 1}
//               onBack={handleBack}
//               onSubmit={handleSubmit}
//               onCreateUpload={handleCreateUpload}
//               onUploadComplete={handleUploadComplete}
//               onUploadDelete={handleUploadDelete}
//               appName={appName}
//             />
//           </Grid>
//         </Grid>
//       </GridContainer>
//     </div>
//   );
// };

// export default Gunsafe;
