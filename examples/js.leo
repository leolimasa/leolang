DomElement = struct
    innerHTML (Fn String Void)
    children (Fn (List DomElement))

Document = struct
    getElementById (Fn String DomElement)
    
extern document Document
extern window Window

render-div = fn ()
    button = (document querySelector) "my-button"

    (button addEventListener) "click"
    	fn ()
	    if not (navigator mediaDevices selectAudioOutput)
	    	(console log) "not supported"
		return Void

	    audio-device v= navigator mediaDevices selectAudioOutput
	    
	    audio = (document createElement) "audio"
	    audio src = "https://example.com"
	    audio play

	    (audio setSinkId) (audioDevice deviceId)
