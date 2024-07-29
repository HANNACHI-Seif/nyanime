package main


func main() {
	anime_name := parse_args();
	instance := ScrapperInstance{};
	instance.Init("nyaa.si", "https://nyaa.si");
	instance.Search(anime_name);
}