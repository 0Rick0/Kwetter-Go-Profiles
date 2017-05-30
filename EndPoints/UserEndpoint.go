package EndPoints

import (
	"log"
	"github.com/emicklei/go-restful"
	"../types"
	"os"
	"io"
	"strings"
	"mime"
)

func (sc *ServiceContainer) DefineUserEndpoints(container *restful.Container)  {
	ws := new(restful.WebService)
	ws.Path("/users").
	Doc("Access to the user resource").
	Consumes(restful.MIME_JSON, restful.MIME_XML).
	Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{username}").
		To(sc.getUser).
		Doc("Get the user's profilepicture").
		Operation("getUser").
		Param(ws.
			PathParameter("username", "The username of the user to get").
			DataType("string")).
		Produces("image/png", "image/jpeg"))

	ws.Route(ws.POST("/{username}").
		To(sc.updateProfilePicture).
		Doc("Update the users profile picture").
		Operation("updateProfilePicture").
		Param(ws.PathParameter("username", "The username of the user").DataType("string")).
		Param(ws.FormParameter("picture", "The profile picture").DataType("file")).
		Consumes("multipart/form-data").Produces("image/png", "image/jpeg"))

	ws.Route(ws.DELETE("/{username}").To(sc.deleteUser).
		Doc("Delete a user from the database, returning the deleted user").
		Operation("deleteUser").
		Param(ws.PathParameter("username", "The username of the user to delete").
		DataType("string")).
		Writes(types.User{}))

	container.Add(ws)
}

func (sc *ServiceContainer) getUser(req *restful.Request, resp *restful.Response) {
	username := req.PathParameter("username")
	log.Printf("Get user %s", username)
	user := sc.Service.GetUserByUsername(username)

	// check if the user is found
	if user.Id == 0{
		//if not, report an service error
		resp.WriteErrorString(404, "User not found")
	}else{
		fo, err := os.OpenFile("./Pictures/" + username, os.O_RDONLY, 0666)
		if err != nil{
			resp.WriteErrorString(400, "Failed to read file")
			return
		}
		resp.AddHeader("content-type", user.MimeType)
		io.Copy(resp, fo)
	}
}

func (sc *ServiceContainer) updateProfilePicture(req *restful.Request, resp *restful.Response)  {
	username := req.PathParameter("username")
	req.Request.ParseMultipartForm(32<<20)
	file, handler, err := req.Request.FormFile("picture")

	if err != nil {
		resp.WriteErrorString(400, "Failed to get image")
		return
	}

	extension := getFileExtension(handler.Filename)
	mimetype := mime.TypeByExtension(extension)

	defer file.Close()
	f, err := os.OpenFile("./Pictures/" + username, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil{
		resp.WriteErrorString(400, "Failed to create file")
		return
	}
	defer f.Close()
	io.Copy(f, file)

	sc.Service.SetProfilePicture(username, mimetype)
	f.Close()
	fo, err := os.OpenFile("./Pictures/" + username, os.O_RDONLY, 066)
	if err != nil{
		resp.WriteErrorString(400, "Failed to read file")
		return
	}
	resp.AddHeader("content-type", mimetype)
	io.Copy(resp, fo)
}

func (sc *ServiceContainer) deleteUser(req *restful.Request, resp *restful.Response)  {
	username := req.PathParameter("username")
	log.Printf("Deleting user %s", username)

	user := sc.Service.GetUserByUsername(username)
	if user.Id == 0 {
		resp.WriteErrorString(404, "User not found")
		return
	}
	if !sc.Service.RemoveUser(*user){
		resp.WriteErrorString(500, "Failed to remove user")
		return
	}
	if err := os.Remove("./Pictures/"+username); err != nil{
		resp.WriteErrorString(500, "Failed to remove file")
		return
	}
	resp.WriteHeader(204)
}

func getFileExtension(filename string) string{
	return filename[strings.LastIndex(filename, "."):]
}
