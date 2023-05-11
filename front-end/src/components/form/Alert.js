// const Alert = props => {
//   return (
//       <div className={"alert " + props.className} role="alert">
//           {props.message}
//       </div>
//   )
// }

// export default Alert
// import React, { useState } from "react";

// const Alert = props => {
//   const [visible, setVisible] = useState(true);

//   const handleClose = () => {
//     setVisible(false);
//   };

//   if (!visible) {
//     return null;
//   }

//   return (
//     <div className={`alert ${props.className} alert-dismissible`} role="alert">
//       <button type="button" className="btn-close" onClick={handleClose} />
//       {props.message}
//     </div>
//   );
// };

// export default Alert;

import React, { useState, useEffect } from "react";

const Alert = props => {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    if (props.message) {
      setVisible(true);

      const timer = setTimeout(() => {
        setVisible(false);
      }, 10000);

      return () => {
        clearTimeout(timer);
      };
    }
  }, [props.message]);

  return (
    <div
    className={`alert ${props.className} ${
      visible ? "slide-in" : "slide-out"
    }`}
      role="alert"
    >
      {props.message}
    </div>
  );
};

export default Alert;